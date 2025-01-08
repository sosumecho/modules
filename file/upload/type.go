package upload

import (
	"sync"
)

// Type 文件类型
type Type struct {
	TypeMap sync.Map
	sync.Once
}

func NewType() *Type {
	return &Type{}
}

/*
	func (t *Type) lazyLoad() {
		t.Do(func() {
			t.TypeMap.Store("ffd8ff", "jpg")                      // JPEG (jpg)
			t.TypeMap.Store("89504e470d0a1a0a0000", "png")        // PNG (png)
			t.TypeMap.Store("474946383961", "gif")                // GIF (gif)
			t.TypeMap.Store("49492a00227105008037", "tif")        // TIFF (tif)
			t.TypeMap.Store("424d228c010000000000", "bmp")        // 16色位图(bmp)
			t.TypeMap.Store("424d8240090000000000", "bmp")        // 24位位图(bmp)
			t.TypeMap.Store("424d8e1b030000000000", "bmp")        // 256色位图(bmp)
			t.TypeMap.Store("41433130313500000000", "dwg")        // CAD (dwg)
			t.TypeMap.Store("3c21444f435459504520", "html")       // HTML (html)   3c68746d6c3e0  3c68746d6c3e0
			t.TypeMap.Store("3c68746d6c3e0", "html")              // HTML (html)   3c68746d6c3e0  3c68746d6c3e0
			t.TypeMap.Store("3c21646f637479706520", "htm")        // HTM (htm)
			t.TypeMap.Store("48544d4c207b0d0a0942", "css")        // css
			t.TypeMap.Store("696b2e71623d696b2e71", "js")         // js
			t.TypeMap.Store("7b5c727466315c616e73", "rtf")        // Rich Text Format (rtf)
			t.TypeMap.Store("38425053000100000000", "psd")        // Photoshop (psd)
			t.TypeMap.Store("46726f6d3a203d3f6762", "eml")        // Email [Outlook Express 6] (eml)
			t.TypeMap.Store("d0cf11e0a1b11ae10000", "doc")        // MS Excel 注意：word、msi 和 excel的文件头一样
			t.TypeMap.Store("d0cf11e0a1b11ae10000", "vsd")        // Visio 绘图
			t.TypeMap.Store("5374616E64617264204A", "mdb")        // MS Access (mdb)
			t.TypeMap.Store("252150532D41646F6265", "ps")         //
			t.TypeMap.Store("255044462d312", "pdf")               // Adobe Acrobat (pdf)
			t.TypeMap.Store("2e524d46000000120001", "rmvb")       // rmvb/rm相同
			t.TypeMap.Store("464c5601050000000900", "flv")        // flv与f4v相同
			t.TypeMap.Store("00000020667479706d70", "mp4")        //
			t.TypeMap.Store("49443303000000002176", "mp3")        //
			t.TypeMap.Store("000001ba210001000180", "mpg")        //
			t.TypeMap.Store("3026b2758e66cf11a6d9", "wmv")        // wmv与asf相同
			t.TypeMap.Store("52494646e27807005741", "wav")        // Wave (wav)
			t.TypeMap.Store("52494646d07d60074156", "avi")        //
			t.TypeMap.Store("4d546864000000060001", "mid")        // MIDI (mid)
			t.TypeMap.Store("504b0304140000000800", "zip")        //
			t.TypeMap.Store("526172211a0700cf9073", "rar")        //
			t.TypeMap.Store("235468697320636f6e66", "ini")        //
			t.TypeMap.Store("504b03040a0000000000", "jar")        //
			t.TypeMap.Store("4d5a9000030000000400", "exe")        // 可执行文件
			t.TypeMap.Store("3c25402070616765206c", "jsp")        // jsp文件
			t.TypeMap.Store("4d616e69666573742d56", "mf")         // MF文件
			t.TypeMap.Store("3c3f786d6c2076657273", "xml")        // xml文件
			t.TypeMap.Store("494e5345525420494e54", "sql")        // xml文件
			t.TypeMap.Store("7061636b616765207765", "java")       // java文件
			t.TypeMap.Store("406563686f206f66660d", "bat")        // bat文件
			t.TypeMap.Store("1f8b0800000000000000", "gz")         // gz文件
			t.TypeMap.Store("6c6f67346a2e726f6f74", "properties") // bat文件
			t.TypeMap.Store("cafebabe0000002e0041", "class")      // bat文件
			t.TypeMap.Store("49545346030000006000", "chm")        // bat文件
			t.TypeMap.Store("04000000010000001300", "mxp")        // bat文件
			t.TypeMap.Store("504b0304140006000800", "docx")       // docx文件
			t.TypeMap.Store("d0cf11e0a1b11ae10000", "wps")        // WPS文字wps、表格et、演示dps都是一样的
			t.TypeMap.Store("6431303a637265617465", "torrent")    //
			t.TypeMap.Store("6D6F6F76", "mov")                    // Quicktime (mov)
			t.TypeMap.Store("FF575043", "wpd")                    // WordPerfect (wpd)
			t.TypeMap.Store("CFAD12FEC5FD746F", "dbx")            // Outlook Express (dbx)
			t.TypeMap.Store("2142444E", "pst")                    // Outlook (pst)
			t.TypeMap.Store("AC9EBD8F", "qdf")                    // Quicken (qdf)
			t.TypeMap.Store("E3828596", "pwl")                    // Windows Password (pwl)
			t.TypeMap.Store("2E7261FD", "ram")                    // Real Audio (ram)
			t.TypeMap.Store("504b0304", "apk")
			t.TypeMap.Store("27444f4d", "txt")
		})
	}

// 获取前面结果字节的二进制

	func (t *Type) bytesToHexString(src []byte) string {
		res := bytes.Buffer{}
		if src == nil || len(src) <= 0 {
			return ""
		}
		temp := make([]byte, 0)
		for _, v := range src {
			sub := v & 0xFF
			hv := hex.EncodeToString(append(temp, sub))
			if len(hv) < 2 {
				res.WriteString(strconv.FormatInt(int64(0), 10))
			}
			res.WriteString(hv)
		}
		return res.String()
	}

// GetType 用文件前面几个字节来判断
// fSrc: 文件字节流（就用前面几个字节）

	func (t *Type) GetType(fSrc []byte) string {
		t.lazyLoad()
		var fileType string
		fileCode := t.bytesToHexString(fSrc)
		t.TypeMap.Range(func(key, value interface{}) bool {
			k := key.(string)
			v := value.(string)
			if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
				strings.HasPrefix(k, strings.ToLower(fileCode)) {
				fileType = v
				return false
			}
			return true
		})
		return fileType
	}
*/
