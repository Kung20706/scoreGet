package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	// 这里假设 bodyBytes 包含 UTF-8 编码的字节序列
	bodyBytes := []byte("�R=o�@��+��R����@(I�4Thm�٫�w����� BG�RHPHB�(q ���o9�r+ Z��B�Y䞲䮪d�YJr[�h)p&�����C��$gڀ ��D=c�L`=�k腞�Sc�� -�J����Z0.}�8����jV��_:E��!I3&�bQA�9��������˱`��� f��:�n��e&���Q���>7�~^�~2{�m�X�8�>>���~�G����h��r6��Ú�����vM~<���h�����θ ���K�LaН^4�خ߼�}�Aŀ���")
	
	// 使用 utf8.DecodeRune 解码字节序列
	for len(bodyBytes) > 0 {
		r, size := utf8.DecodeRune(bodyBytes)
		if r == utf8.RuneError && size == 1 {
			fmt.Println("Error decoding rune")
			break
		}
		// 处理解码后的 Unicode 字符 r
		fmt.Printf("%c", r)
		// 调整切片，准备解码下一个字符
		bodyBytes = bodyBytes[size:]
	}
}

