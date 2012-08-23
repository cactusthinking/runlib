package win32

// +build windows,386

import (
	"syscall"
	"unsafe"
	"runlib/platform"
)

var (
	procGetUserObjectSecurity = user32.NewProc("GetUserObjectSecurity")
	procGetSecurityDescriptorDacl = advapi32.NewProc("GetSecurityDescriptorDacl")
	procIsValidAcl = advapi32.NewProc("IsValidAcl")
)

const (
	DACL_SECURITY_INFORMATION = 0x00000004
)

func AddAceToDesktop(desk Hdesk, sid *syscall.SID) {
	secDesc, err := GetUserObjectSecurity(syscall.Handle(desk), DACL_SECURITY_INFORMATION)
}	


func CreateSecurityDescriptor(length int) []byte {
	return platform.AlignedBuffer(length, 4)
}
	

func GetUserObjectSecurity_Ex(obj syscall.Handle, sid uint32, desc []byte) (uint32, error) {
	var nLength uint32
	r1, _, e1 := procGetUserObjectSecurity.Call(
		uintptr(obj),
		uintptr(unsafe.Pointer(&sid)),
		uintptr(unsafe.Pointer(&desc[0])),
		uintptr(len(desc)),
		uintptr(unsafe.Pointer(&nLength)))
	if int(r1) == 0 {
		return nLength, e1
	}
	return nLength, nil
}

func GetUserObjectSecurity(obj syscall.Handle, sid uint32) ([]byte, error) {
	nLength, err := GetUserObjectSecurity_Ex(obj, sid, nil)
	if nLength == 0 {
		return nil, err
	}

	desc := CreateSecurityDescriptor(int(nLength))
	_, err = GetUserObjectSecurity_Ex(obj, sid, desc)
	if err != nil {
		return nil, err
	}
	return desc, err
}

type Acl struct{}

func GetSecurityDescriptorDacl(sid []byte) (present bool, acl *Acl, defaulted bool, err error) {
	var dPresent, dDefaulted uint32
	r1, _, e1 := procGetSecurityDescriptorDacl.Call(
		uintptr(unsafe.Pointer(&sid[0])),
		uintptr(unsafe.Pointer(&dPresent)),
		uintptr(unsafe.Pointer(&acl)),
		uintptr(unsafe.Pointer(&dDefaulted)))
	if dPresent != 0 {
		present = true
	}
	if dDefaulted != 0 {
		defaulted = true
	}
	if int(r1) == 0 {
		err = e1
	}
	return
}

func IsValidAcl(acl *Acl) bool {
	r1, _, _ := procIsValidAcl.Call(
		uintptr(unsafe.Pointer(acl)))
	if r1 == 0 {
		return false
	}
	return true
}