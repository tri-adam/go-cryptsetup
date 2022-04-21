package cryptsetup

import (
	"testing"
)

func Test_LUKS1_Format(test *testing.T) {
	testWrapper := TestWrapper{test}

	device, err := Init(DevicePath)
	testWrapper.AssertNoError(err)

	hashBeforeFormat := getFileMD5(DevicePath, test)

	err = device.Format(LUKS1{Hash: "sha256"}, GenericParams{Cipher: "aes", CipherMode: "xts-plain64", VolumeKeySize: 512 / 8})
	testWrapper.AssertNoError(err)

	hashAfterFormat := getFileMD5(DevicePath, test)

	if hashBeforeFormat == hashAfterFormat {
		test.Error("Unsuccessful call to Format() when using LUKS1 parameters.")
	}

	if device.Type() != "LUKS1" {
		test.Error("Expected type: LUKS1.")
	}

	device.Free()
}

func Test_LUKS1_Load_ActivateByPassphrase_Deactivate(test *testing.T) {
	testWrapper := TestWrapper{test}
	luks1 := LUKS1{Hash: "sha256"}

	device, err := Init(DevicePath)
	testWrapper.AssertNoError(err)
	err = device.Format(luks1, GenericParams{Cipher: "aes", CipherMode: "xts-plain64", VolumeKeySize: 512 / 8})
	testWrapper.AssertNoError(err)

	err = device.KeyslotAddByVolumeKey(0, "", "testPassphrase")
	testWrapper.AssertNoError(err)

	device.Free()

	device, err = Init(DevicePath)
	testWrapper.AssertNoError(err)
	err = device.Load(nil)
	testWrapper.AssertNoError(err)

	err = device.ActivateByPassphrase(DeviceName, 0, "testPassphrase", CRYPT_ACTIVATE_READONLY)
	testWrapper.AssertNoError(err)

	err = device.Deactivate(DeviceName)
	testWrapper.AssertNoError(err)

	if device.Type() != "LUKS1" {
		test.Error("Expected type: LUKS1.")
	}

	device.Free()
}

func Test_LUKS1_Load_ActivateByPassphrase_Free_InitByName_Deactivate(test *testing.T) {
	testWrapper := TestWrapper{test}
	luks1 := LUKS1{Hash: "sha256"}

	device, err := Init(DevicePath)
	testWrapper.AssertNoError(err)
	err = device.Format(luks1, GenericParams{Cipher: "aes", CipherMode: "xts-plain64", VolumeKeySize: 512 / 8})
	testWrapper.AssertNoError(err)

	err = device.KeyslotAddByVolumeKey(0, "", "testPassphrase")
	testWrapper.AssertNoError(err)

	device.Free()

	device, err = Init(DevicePath)
	testWrapper.AssertNoError(err)
	err = device.Load(nil)
	testWrapper.AssertNoError(err)

	err = device.ActivateByPassphrase(DeviceName, 0, "testPassphrase", CRYPT_ACTIVATE_READONLY)
	testWrapper.AssertNoError(err)

	device.Free()

	device, err = InitByName(DeviceName)
	testWrapper.AssertNoError(err)

	err = device.Deactivate(DeviceName)
	testWrapper.AssertNoError(err)

	if device.Type() != "LUKS1" {
		test.Error("Expected type: LUKS1.")
	}

	device.Free()
}

func Test_LUKS1_ActivateByVolumeKey_Deactivate(test *testing.T) {
	testWrapper := TestWrapper{test}

	genericParams := GenericParams{
		Cipher:        "aes",
		CipherMode:    "xts-plain64",
		VolumeKey:     generateKey(512/8, test),
		VolumeKeySize: 512 / 8,
	}

	device, err := Init(DevicePath)
	testWrapper.AssertNoError(err)

	err = device.Format(LUKS1{Hash: "sha256"}, genericParams)
	testWrapper.AssertNoError(err)

	err = device.ActivateByVolumeKey(DeviceName, genericParams.VolumeKey, genericParams.VolumeKeySize, CRYPT_ACTIVATE_READONLY)
	testWrapper.AssertNoError(err)

	err = device.Deactivate(DeviceName)
	testWrapper.AssertNoError(err)

	if device.Type() != "LUKS1" {
		test.Error("Expected type: LUKS1.")
	}

	device.Free()
}

func Test_LUKS1_ActivateByAutoGeneratedVolumeKey_Deactivate(test *testing.T) {
	testWrapper := TestWrapper{test}

	genericParams := GenericParams{
		Cipher:        "aes",
		CipherMode:    "xts-plain64",
		VolumeKeySize: 512 / 8,
	}

	device, err := Init(DevicePath)
	testWrapper.AssertNoError(err)

	err = device.Format(LUKS1{Hash: "sha256"}, genericParams)
	testWrapper.AssertNoError(err)

	err = device.ActivateByVolumeKey(DeviceName, "", 512/8, CRYPT_ACTIVATE_READONLY)
	testWrapper.AssertNoError(err)

	err = device.Deactivate(DeviceName)
	testWrapper.AssertNoError(err)

	if device.Type() != "LUKS1" {
		test.Error("Expected type: LUKS1.")
	}

	device.Free()
}

func Test_LUKS1_KeyslotAddByVolumeKey(test *testing.T) {
	testWrapper := TestWrapper{test}

	device, err := Init(DevicePath)
	testWrapper.AssertNoError(err)

	err = device.Format(LUKS1{Hash: "sha256"}, GenericParams{Cipher: "aes", CipherMode: "xts-plain64", VolumeKeySize: 512 / 8})
	testWrapper.AssertNoError(err)

	err = device.KeyslotAddByVolumeKey(0, "", "testPassphrase")
	testWrapper.AssertNoError(err)

	err = device.KeyslotAddByVolumeKey(0, "", "testPassphrase")
	testWrapper.AssertError(err)
	testWrapper.AssertErrorCodeEquals(err, -22)

	device.Free()
}

func Test_LUKS1_KeyslotAddByPassphrase(test *testing.T) {
	testWrapper := TestWrapper{test}

	device, err := Init(DevicePath)
	testWrapper.AssertNoError(err)

	err = device.Format(LUKS1{Hash: "sha256"}, GenericParams{Cipher: "aes", CipherMode: "xts-plain64", VolumeKeySize: 512 / 8})
	testWrapper.AssertNoError(err)

	err = device.KeyslotAddByVolumeKey(0, "", "testPassphrase")
	testWrapper.AssertNoError(err)

	err = device.KeyslotAddByPassphrase(1, "testPassphrase", "secondTestPassphrase")
	testWrapper.AssertNoError(err)

	err = device.KeyslotAddByPassphrase(1, "testPassphrase", "secondTestPassphrase")
	testWrapper.AssertError(err)
	testWrapper.AssertErrorCodeEquals(err, -22)

	device.Free()
}

func Test_LUKS1_KeyslotChangeByPassphrase(test *testing.T) {
	testWrapper := TestWrapper{test}

	device, err := Init(DevicePath)
	testWrapper.AssertNoError(err)

	err = device.Format(LUKS1{Hash: "sha256"}, GenericParams{Cipher: "aes", CipherMode: "xts-plain64", VolumeKeySize: 512 / 8})
	testWrapper.AssertNoError(err)

	err = device.KeyslotAddByVolumeKey(0, "", "testPassphrase")
	testWrapper.AssertNoError(err)

	err = device.KeyslotChangeByPassphrase(0, 0, "testPassphrase", "secondTestPassphrase")
	testWrapper.AssertNoError(err)

	err = device.ActivateByPassphrase(DeviceName, 0, "secondTestPassphrase", CRYPT_ACTIVATE_READONLY)
	testWrapper.AssertNoError(err)

	err = device.Deactivate(DeviceName)
	testWrapper.AssertNoError(err)

	err = device.ActivateByPassphrase(DeviceName, 0, "testPassphrase", CRYPT_ACTIVATE_READONLY)
	testWrapper.AssertError(err)
	testWrapper.AssertErrorCodeEquals(err, -1)

	device.Free()
}
