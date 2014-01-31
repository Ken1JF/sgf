package sgf_test

import (
	"fmt"
	"github.com/Ken1JF/ah"
	"github.com/Ken1JF/sgf"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const SGFDir = "./testdata"
const OutDir = "./testout"

// count should be 56 with all tests included.
// One test currently in "WorkOnLater"
const outputVerified int = 54      // controls the printing of last graph
const testStringsVerified int = 55 // controls tracing

const SGFDB_NUM_PER_LINE = 12

func ExampleReadWriteSGFFile() {
	err := sgf.SetupSGFProperties(defaultSpecFile, false, false)
	if err == 0 {
		fmt.Println("Reading SGF Tests, SGFDir =", SGFDir, ", outDir =", OutDir)
		dirFiles, err := ioutil.ReadDir(SGFDir)
		if err != nil && err != io.EOF {
			fmt.Println("Can't read test directory:", SGFDir)
			return
		}
		// Check the output directory. If missing, create it.
		_, errS := os.Stat(OutDir)
		if errS != nil {
			err2 := os.MkdirAll(OutDir, os.ModeDir|os.ModePerm)
			if err2 != nil {
				fmt.Println("ReadAndWriteDirectory Error:", err2, "trying to create test output directory:", OutDir)
				fmt.Println("Original Error:", err, "trying os.Stat")
				return
			}
		}

		count := 0
		for _, f := range dirFiles {
			if strings.Index(f.Name(), ".sgf") >= 0 {
				count += 1
				fmt.Println("Processing:", f.Name())
				if count > testStringsVerified {
					ah.SetAHTrace(true)
					fmt.Println("Tracing", f.Name())
				}
				fileName := SGFDir + "/" + f.Name()
				b, err := ioutil.ReadFile(fileName)
				if err != nil && err != io.EOF {
					fmt.Println("Error reading file:", fileName, err)
					return
				}
				//			prsr , errL := sgf.ParseFile(fileName, b, sgf.ParseComments, 0)
				prsr, errL := sgf.ParseFile(fileName, b, sgf.ParseComments+sgf.ParserPlay, 0)
				if len(errL) != 0 {
					fmt.Println("Error while parsing:", fileName, ", ", errL.Error())
					return
				}
				outFileName := OutDir + "/" + f.Name()
				err = prsr.GameTree.WriteFile(outFileName, SGFDB_NUM_PER_LINE)
				if err != nil {
					fmt.Printf("Error writing: %s, %s\n", outFileName, err)
				}
				if count > outputVerified {
					prsr.PrintAbstHier(fileName, true)
				}
			}
		}
	}

	// Output:
	// Reading SGF Tests, SGFDir = ./testdata , outDir = ./testout
	// Processing: 0.sgf
	// Processing: 0A.sgf
	// Processing: 0B.sgf
	// Processing: 0C.sgf
	// Processing: 0D.sgf
	// Processing: 0E.sgf
	// Processing: 0F.sgf
	// Processing: 0G.sgf
	// Processing: 0H.sgf
	// Processing: 0K.sgf
	// Processing: 0L.sgf
	// Processing: 0M.sgf
	// Processing: 1.sgf
	// Processing: 10.sgf
	// ./testdata/10.sgf:1:33: Unknown SGF property: WHAT:bb
	// Processing: 1000.sgf
	// Processing: 1001.sgf
	// Processing: 1002.sgf
	// Processing: 1003.sgf
	// Processing: 1004.sgf
	// Processing: 1005.sgf
	// Processing: 1006.sgf
	// Processing: 1007.sgf
	// Processing: 1008.sgf
	// Processing: 1009.sgf
	// Processing: 1010.sgf
	// Processing: 1011.sgf
	// Processing: 1012.sgf
	// Processing: 1013.sgf
	// Processing: 1014.sgf
	// Processing: 1015.sgf
	// Processing: 1016.sgf
	// Processing: 1017.sgf
	// Processing: 1018.sgf
	// Processing: 1019.sgf
	// Processing: 1020.sgf
	// Processing: 1021.sgf
	// Processing: 1022.sgf
	// Processing: 1023.sgf
	// Processing: 1024.sgf
	// Processing: 1025.sgf
	// Processing: 2.sgf
	// Processing: 3.sgf
	// Processing: 4.sgf
	// Processing: 5.sgf
	// Processing: 6.sgf
	// Processing: Lyonweiqi-ken1jf-2.sgf
	// Processing: Lyonweiqi-ken1jf-Apr_26.sgf
	// Processing: Lyonweiqi-ken1jf_4.sgf
	// Processing: Lyonweiqi-ken1jf_May_30.sgf
	// Processing: Lyonweiqi-ken1jf_May_7_2008.sgf
	// Processing: Lyonweiqi-ken1jf_May_9_2008.sgf
	// Processing: Lyonweiqi-scottSsgf.sgf
	// Processing: Lyonweiqi.sgf
	// Processing: print1.sgf
	// Processing: print2.sgf
	// Abstraction Hierarchy: ./testdata/print2.sgf
	// Level 1
	// Black nodes
	// 57:1,1 1-mem:(Q,16):1:1,adj:270(1),269(1),164(1),252(1),
	// 66:1,1 4-mem:(S,4):1:1,(P,4),(Q,4),(R,4),adj:19(1),12(3),285(1),281(1),272(1),105(1),205(2),
	// 67:1,1 1-mem:(F,3):1:1,adj:56(1),6(1),180(1),177(1),
	// 71:1,1 1-mem:(Q,6):1:1,adj:59(1),286(1),204(1),205(1),
	// 82:1,1 2-mem:(J,3):1:1,(J,4),adj:79(1),74(1),73(1),16(1),275(1),92(1),
	// 88:1,1 1-mem:(L,1):1:1,adj:276(1),91(1),185(1),
	// 93:1,1 2-mem:(O,3):1:1,(N,3),adj:19(1),12(1),83(1),279(1),278(1),183(1),
	// 111:1,1 1-mem:(M,17):1:1,adj:265(1),267(1),258(1),266(1),
	// 186:1,1 1-mem:(K,2):1:1,adj:92(1),185(1),276(1),275(1),
	// 201:1,1 1-mem:(Q,10):1:1,adj:20(1),43(1),35(1),203(1),
	// 280:1,1 1-mem:(M,2):1:1,adj:279(1),276(1),83(1),91(1),
	// Total 11 nodes, with 16 members
	// White nodes
	// 12:2,2 3-mem:(P,3):2:2,(Q,3),(R,3),adj:273(1),272(1),190(1),93(1),66(3),85(1),
	// 19:2,2 1-mem:(O,4):2:2,adj:283(1),66(1),278(1),93(1),
	// 25:2,2 1-mem:(C,14):2:2,adj:9(1),166(1),47(1),78(1),
	// 55:2,2 1-mem:(K,16):2:2,adj:37(1),54(1),261(1),168(1),
	// 83:2,2 1-mem:(M,3):2:2,adj:18(1),86(1),280(1),93(1),
	// 98:2,2 1-mem:(D,17):2:2,adj:53(1),62(1),58(1),146(1),
	// 181:2,2 1-mem:(D,4):2:2,adj:44(1),271(1),196(1),195(1),
	// 193:2,2 1-mem:(D,6):2:2,adj:132(1),192(1),209(1),196(1),
	// 205:2,2 2-mem:(Q,5):2:2,(P,5),adj:71(1),204(1),66(2),285(1),283(1),
	// 274:2,2 2-mem:(L,5):2:2,(L,4),adj:86(1),18(1),81(1),17(1),23(1),79(1),
	// 275:2,2 1-mem:(K,3):2:2,adj:86(1),79(1),82(1),186(1),
	// 276:2,2 1-mem:(L,2):2:2,adj:86(1),186(1),88(1),280(1),
	// Total 12 nodes, with 16 members
	// Unocc nodes
	// 0:202,3 1-mem:(T,19):202:202,adj:65(1),8(1),
	// 1:778,3 1-mem:(A,1):778:778,adj:69(1),2(1),
	// 2:906,3 8-mem:(J,1):906:906,(H,1),(G,1),(F,1),(E,1),(D,1),(C,1),(B,1),adj:1(1),185(1),99(1),5(1),133(1),100(1),92(2),6(2),
	// 3:394,3 1-mem:(T,1):394:394,adj:188(1),31(1),
	// 4:842,3 7-mem:(A,18):842:842,(A,17),(A,16),(A,15),(A,14),(A,13),(A,12),adj:110(1),64(1),61(1),9(6),
	// 5:3018,3 1-mem:(B,2):3018:3018,adj:2(1),100(1),102(1),69(1),
	// 6:4042,3 2-mem:(E,2):4042:4042,(F,2),adj:2(2),177(1),99(1),133(1),67(1),
	// 7:12234,3 1-mem:(K,6):12234:12234,adj:17(1),154(1),33(1),23(1),
	// 8:458,3 12-mem:(T,18):458:458,(T,17),(T,16),(T,15),(T,14),(T,13),(T,12),(T,11),(T,10),(T,9),(T,8),(T,7),adj:0(1),247(1),84(1),34(1),13(10),
	// 9:4042,3 6-mem:(B,17):4042:4042,(B,16),(B,15),(B,14),(B,13),(B,12),adj:4(6),225(1),166(1),151(1),145(1),25(1),78(1),61(1),58(1),
	// 10:24786,3 1-mem:(D,13):24786:24786,adj:141(1),224(1),47(1),78(1),
	// 11:5259,3 1-mem:(T,3):5259:5259,adj:272(1),105(1),188(1),
	// 13:4042,3 10-mem:(S,17):4042:4042,(S,16),(S,15),(S,14),(S,13),(S,12),(S,11),(S,10),(S,9),(S,8),adj:8(10),277(1),270(1),247(1),115(1),84(1),20(7),
	// 14:17163,3 1-mem:(N,9):17163:17163,adj:241(1),107(1),199(1),126(1),
	// 15:18509,3 1-mem:(B,10):18509:18509,adj:178(1),120(1),151(1),118(1),
	// 16:8138,3 1-mem:(H,4):8138:8138,adj:114(1),94(1),82(1),73(1),
	// 17:10900,3 1-mem:(L,6):10900:10900,adj:7(1),274(1),113(1),28(1),
	// 18:6867,3 1-mem:(M,4):6867:6867,adj:278(1),81(1),274(1),83(1),
	// 20:6090,3 7-mem:(R,15):6090:6090,(R,14),(R,13),(R,12),(R,11),(R,10),(R,9),adj:13(7),277(1),270(1),164(1),201(1),35(1),43(4),
	// 21:14540,3 1-mem:(D,8):14540:14540,adj:125(1),96(1),209(1),135(1),
	// 22:16403,3 1-mem:(A,9):16403:16403,adj:120(1),150(1),118(1),
	// 23:8781,3 1-mem:(K,5):8781:8781,adj:7(1),274(1),74(1),79(1),
	// 24:27101,3 1-mem:(H,14):27101:27101,adj:52(1),158(1),49(1),106(1),
	// 26:10186,3 1-mem:(E,14):10186:10186,adj:141(1),159(1),47(1),50(1),
	// 27:6090,3 1-mem:(P,17):6090:6090,adj:252(1),60(1),264(1),269(1),
	// 28:12234,3 2-mem:(N,6):12234:12234,(M,6),adj:17(1),284(1),29(1),240(1),128(1),81(1),
	// 29:11210,3 1-mem:(O,6):11210:11210,adj:28(1),283(1),204(1),243(1),
	// 30:10186,3 1-mem:(P,13):10186:10186,adj:162(1),242(1),160(1),43(1),
	// 31:1104,3 1-mem:(S,1):1104:1104,adj:3(1),189(1),89(1),
	// 32:13258,3 1-mem:(G,7):13258:13258,adj:36(1),220(1),130(1),39(1),
	// 33:10770,3 1-mem:(J,6):10770:10770,adj:7(1),97(1),42(1),74(1),
	// 34:11408,3 1-mem:(T,6):11408:11408,adj:8(1),282(1),139(1),
	// 35:17362,3 1-mem:(Q,9):17362:17362,adj:20(1),202(1),201(1),200(1),
	// 36:14282,3 4-mem:(G,11):14282:14282,(G,8),(G,9),(G,10),adj:32(1),236(1),237(1),238(1),222(1),149(1),227(1),148(1),234(1),147(1),
	// 37:29259,3 1-mem:(K,15):29259:29259,adj:87(1),55(1),51(1),48(1),
	// 38:14282,3 1-mem:(N,12):14282:14282,adj:46(1),219(1),161(1),155(1),
	// 39:10638,3 1-mem:(G,6):10638:10638,adj:32(1),208(1),42(1),121(1),
	// 40:1994,3 1-mem:(K,10):1994:1994,adj:231(1),134(1),228(1),232(1),
	// 41:15179,3 1-mem:(O,8):15179:15179,adj:243(1),206(1),107(1),241(1),
	// 42:10699,3 1-mem:(H,6):10699:10699,adj:33(1),39(1),220(1),114(1),
	// 43:8138,3 4-mem:(Q,14):8138:8138,(Q,13),(Q,12),(Q,11),adj:20(4),30(1),162(1),164(1),201(1),170(1),160(1),
	// 44:6290,3 1-mem:(C,4):6290:6290,adj:75(1),194(1),181(1),102(1),
	// 45:16666,3 1-mem:(E,9):16666:16666,adj:96(1),137(1),135(1),148(1),
	// 46:13258,3 1-mem:(N,13):13258:13258,adj:38(1),163(1),242(1),216(1),
	// 47:26830,3 1-mem:(D,14):26830:26830,adj:26(1),25(1),10(1),70(1),
	// 48:27230,3 1-mem:(K,14):27230:27230,adj:37(1),138(1),63(1),158(1),
	// 49:29132,3 1-mem:(H,15):29132:29132,adj:24(1),144(1),51(1),54(1),
	// 50:9162,3 1-mem:(E,15):9162:9162,adj:26(1),142(1),70(1),123(1),
	// 51:29204,3 1-mem:(J,15):29204:29204,adj:49(1),37(1),158(1),54(1),
	// 52:12234,3 1-mem:(G,14):12234:12234,adj:24(1),144(1),159(1),143(1),
	// 53:33043,3 1-mem:(E,17):33043:33043,adj:253(1),123(1),98(1),62(1),
	// 54:8138,3 2-mem:(J,16):8138:8138,(H,16),adj:49(1),51(1),117(1),248(1),131(1),55(1),
	// 56:4492,3 1-mem:(G,3):4492:4492,adj:73(1),94(1),67(1),133(1),
	// 58:5066,3 1-mem:(C,17):5066:5066,adj:9(1),145(1),98(1),62(1),
	// 59:13264,3 1-mem:(Q,7):13264:13264,adj:80(1),202(1),207(1),71(1),
	// 60:35723,3 1-mem:(P,18):35723:35723,adj:27(1),251(1),174(1),260(1),
	// 61:3018,3 1-mem:(B,18):3018:3018,adj:4(1),9(1),116(1),62(1),
	// 62:4042,3 3-mem:(E,18):4042:4042,(D,18),(C,18),adj:53(1),61(1),58(1),175(1),98(1),116(3),
	// 63:27280,3 1-mem:(L,14):27280:27280,adj:48(1),218(1),217(1),87(1),
	// 64:586,3 1-mem:(A,19):586:586,adj:4(1),116(1),
	// 65:37968,3 1-mem:(S,19):37968:37968,adj:0(1),247(1),101(1),
	// 68:37328,3 1-mem:(H,19):37328:37328,adj:76(1),77(1),176(1),
	// 69:842,3 3-mem:(A,2):842:842,(A,3),(A,4),adj:5(1),1(1),72(1),102(2),
	// 70:8138,3 1-mem:(D,15):8138:8138,adj:50(1),47(1),146(1),166(1),
	// 72:8208,3 1-mem:(A,5):8208:8208,adj:69(1),182(1),90(1),
	// 73:6090,3 1-mem:(H,3):6090:6090,adj:56(1),16(1),92(1),82(1),
	// 74:8721,3 1-mem:(J,5):8721:8721,adj:23(1),33(1),114(1),82(1),
	// 75:5066,3 1-mem:(C,3):5066:5066,adj:44(1),271(1),102(1),100(1),
	// 76:35275,3 1-mem:(H,18):35275:35275,adj:68(1),257(1),254(1),248(1),
	// 77:37387,3 1-mem:(J,19):37387:37387,adj:68(1),257(1),255(1),
	// 78:24723,3 1-mem:(C,13):24723:24723,adj:10(1),25(1),9(1),225(1),
	// 79:6748,3 1-mem:(K,4):6748:6748,adj:23(1),275(1),274(1),82(1),
	// 80:13323,3 1-mem:(R,7):13323:13323,adj:59(1),286(1),84(1),277(1),
	// 81:10186,3 1-mem:(M,5):10186:10186,adj:18(1),28(1),284(1),274(1),
	// 84:13392,3 1-mem:(S,7):13392:13392,adj:80(1),8(1),13(1),282(1),
	// 85:4042,3 1-mem:(Q,2):4042:4042,adj:12(1),273(1),190(1),191(1),
	// 86:4756,3 1-mem:(L,3):4756:4756,adj:83(1),276(1),274(1),275(1),
	// 87:29332,3 1-mem:(L,15):29332:29332,adj:63(1),37(1),169(1),168(1),
	// 89:1038,3 1-mem:(R,1):1038:1038,adj:31(1),190(1),191(1),
	// 90:10251,3 1-mem:(A,6):10251:10251,adj:72(1),245(1),152(1),
	// 91:716,3 1-mem:(M,1):716:716,adj:88(1),280(1),187(1),
	// 92:4042,3 2-mem:(H,2):4042:4042,(J,2),adj:73(1),82(1),2(2),133(1),186(1),
	// 94:6541,3 1-mem:(G,4):6541:6541,adj:56(1),16(1),180(1),121(1),
	// 95:18638,3 1-mem:(D,10):18638:18638,adj:178(1),137(1),226(1),135(1),
	// 96:10186,3 1-mem:(E,8):10186:10186,adj:45(1),21(1),149(1),210(1),
	// 97:12814,3 1-mem:(J,7):12814:12814,adj:33(1),221(1),220(1),154(1),
	// 99:2251,3 1-mem:(D,2):2251:2251,adj:6(1),2(1),271(1),100(1),
	// 100:4042,3 1-mem:(C,2):4042:4042,adj:75(1),99(1),5(1),2(1),
	// 101:714,3 1-mem:(R,19):714:714,adj:65(1),250(1),173(1),
	// 102:4042,3 2-mem:(B,3):4042:4042,(B,4),adj:75(1),5(1),44(1),69(2),182(1),
	// 103:14412,3 1-mem:(B,8):14412:14412,adj:125(1),120(1),150(1),244(1),
	// 104:29579,3 1-mem:(P,15):29579:29579,adj:269(1),164(1),165(1),162(1),
	// 105:7312,3 1-mem:(T,4):7312:7312,adj:11(1),66(1),139(1),
	// 106:25043,3 1-mem:(H,13):25043:25043,adj:24(1),235(1),215(1),143(1),
	// 107:17232,3 1-mem:(O,9):17232:17232,adj:41(1),14(1),200(1),157(1),
	// 108:8138,3 1-mem:(F,16):8138:8138,adj:253(1),142(1),123(1),131(1),
	// 109:27476,3 1-mem:(O,14):27476:27476,adj:165(1),162(1),163(1),242(1),
	// 110:20493,3 1-mem:(A,11):20493:20493,adj:4(1),151(1),118(1),
	// 112:16526,3 1-mem:(C,9):16526:16526,adj:125(1),178(1),135(1),120(1),
	// 113:12940,3 1-mem:(L,7):12940:12940,adj:17(1),154(1),128(1),129(1),
	// 114:8652,3 1-mem:(H,5):8652:8652,adj:42(1),74(1),16(1),121(1),
	// 115:33818,3 1-mem:(R,17):33818:33818,adj:13(1),270(1),250(1),252(1),
	// 116:714,3 4-mem:(B,19):714:714,(C,19),(D,19),(E,19),adj:64(1),61(1),62(3),167(1),
	// 117:33293,3 1-mem:(J,17):33293:33293,adj:54(1),261(1),257(1),248(1),
	// 118:842,3 1-mem:(A,10):842:842,adj:15(1),110(1),22(1),
	// 119:714,3 2-mem:(L,19):714:714,(M,19),adj:258(1),255(1),136(1),171(1),
	// 120:16462,3 1-mem:(B,9):16462:16462,adj:103(1),112(1),15(1),22(1),
	// 121:8604,3 1-mem:(G,5):8604:8604,adj:39(1),114(1),94(1),197(1),
	// 122:970,3 1-mem:(L,12):970:970,adj:219(1),218(1),214(1),233(1),
	// 123:30988,3 1-mem:(E,16):30988:30988,adj:108(1),53(1),50(1),146(1),
	// 124:22796,3 1-mem:(E,12):22796:22796,adj:229(1),213(1),141(1),224(1),
	// 125:6090,3 1-mem:(C,8):6090:6090,adj:21(1),112(1),103(1),246(1),
	// 126:17100,3 1-mem:(M,9):17100:17100,adj:14(1),230(1),239(1),127(1),
	// 127:15053,3 1-mem:(M,8):15053:15053,adj:126(1),241(1),129(1),128(1),
	// 128:13004,3 1-mem:(M,7):13004:13004,adj:127(1),113(1),28(1),240(1),
	// 129:14988,3 1-mem:(L,8):14988:14988,adj:127(1),113(1),239(1),153(1),
	// 130:12234,3 1-mem:(F,7):12234:12234,adj:32(1),149(1),210(1),208(1),
	// 131:31118,3 1-mem:(G,16):31118:31118,adj:54(1),108(1),249(1),144(1),
	// 132:10507,3 1-mem:(E,6):10507:10507,adj:208(1),210(1),193(1),198(1),
	// 133:2446,3 1-mem:(G,2):2446:2446,adj:92(1),56(1),6(1),2(1),
	// 134:970,3 2-mem:(J,11):970:970,(K,11),adj:40(1),212(1),236(1),228(1),214(1),233(1),
	// 135:16592,3 1-mem:(D,9):16592:16592,adj:112(1),95(1),45(1),21(1),
	// 136:35474,3 1-mem:(L,18):35474:35474,adj:119(1),258(1),266(1),256(1),
	// 137:10186,3 1-mem:(E,10):10186:10186,adj:95(1),45(1),147(1),229(1),
	// 138:25168,3 1-mem:(K,13):25168:25168,adj:48(1),214(1),218(1),215(1),
	// 139:9356,3 1-mem:(T,5):9356:9356,adj:105(1),34(1),281(1),
	// 140:24922,3 1-mem:(F,13):24922:24922,adj:213(1),141(1),159(1),143(1),
	// 141:24843,3 1-mem:(E,13):24843:24843,adj:124(1),140(1),26(1),10(1),
	// 142:10186,3 1-mem:(F,15):10186:10186,adj:108(1),50(1),144(1),159(1),
	// 143:24974,3 1-mem:(G,13):24974:24974,adj:52(1),106(1),140(1),237(1),
	// 144:29070,3 1-mem:(G,15):29070:29070,adj:142(1),52(1),49(1),131(1),
	// 145:30862,3 1-mem:(C,16):30862:30862,adj:58(1),9(1),166(1),146(1),
	// 146:30926,3 1-mem:(D,16):30926:30926,adj:70(1),123(1),98(1),145(1),
	// 147:18771,3 1-mem:(F,10):18771:18771,adj:137(1),36(1),238(1),148(1),
	// 148:16716,3 1-mem:(F,9):16716:16716,adj:36(1),147(1),45(1),149(1),
	// 149:14675,3 1-mem:(F,8):14675:14675,adj:96(1),36(1),130(1),148(1),
	// 150:14350,3 1-mem:(A,8):14350:14350,adj:103(1),22(1),152(1),
	// 151:20558,3 1-mem:(B,11):20558:20558,adj:15(1),9(1),110(1),179(1),
	// 152:12302,3 1-mem:(A,7):12302:12302,adj:90(1),150(1),244(1),
	// 153:14930,3 1-mem:(K,8):14930:14930,adj:129(1),232(1),221(1),154(1),
	// 154:12878,3 1-mem:(K,7):12878:12878,adj:97(1),7(1),113(1),153(1),
	// 155:21267,3 1-mem:(N,11):21267:21267,adj:38(1),211(1),199(1),156(1),
	// 156:12234,3 1-mem:(O,11):12234:12234,adj:155(1),170(1),161(1),157(1),
	// 157:19293,3 1-mem:(O,10):19293:19293,adj:107(1),156(1),203(1),199(1),
	// 158:27147,3 1-mem:(J,14):27147:27147,adj:48(1),51(1),24(1),215(1),
	// 159:26958,3 1-mem:(F,14):26958:26958,adj:142(1),52(1),26(1),140(1),
	// 160:23438,3 1-mem:(P,12):23438:23438,adj:43(1),30(1),170(1),161(1),
	// 161:23374,3 1-mem:(O,12):23374:23374,adj:156(1),160(1),38(1),242(1),
	// 162:27532,3 1-mem:(P,14):27532:27532,adj:104(1),43(1),109(1),30(1),
	// 163:27404,3 1-mem:(N,14):27404:27404,adj:109(1),46(1),262(1),217(1),
	// 164:29648,3 1-mem:(Q,15):29648:29648,adj:104(1),20(1),57(1),43(1),
	// 165:29516,3 1-mem:(O,15):29516:29516,adj:104(1),109(1),263(1),262(1),
	// 166:28814,3 1-mem:(C,15):28814:28814,adj:70(1),145(1),9(1),25(1),
	// 167:37198,3 1-mem:(F,19):37198:37198,adj:116(1),176(1),175(1),
	// 168:31379,3 1-mem:(L,16):31379:31379,adj:55(1),87(1),267(1),266(1),
	// 169:10186,3 1-mem:(M,15):10186:10186,adj:87(1),262(1),267(1),217(1),
	// 170:21390,3 1-mem:(P,11):21390:21390,adj:156(1),43(1),160(1),203(1),
	// 171:37648,3 1-mem:(N,19):37648:37648,adj:119(1),259(1),172(1),
	// 172:37707,3 1-mem:(O,19):37707:37707,adj:171(1),260(1),174(1),
	// 173:37838,3 1-mem:(Q,19):37838:37838,adj:101(1),251(1),174(1),
	// 174:37774,3 1-mem:(P,19):37774:37774,adj:60(1),172(1),173(1),
	// 175:35150,3 1-mem:(F,18):35150:35150,adj:167(1),62(1),254(1),253(1),
	// 176:37260,3 1-mem:(G,19):37260:37260,adj:68(1),167(1),254(1),
	// 177:4371,3 1-mem:(E,3):4371:4371,adj:6(1),67(1),271(1),195(1),
	// 178:6090,3 1-mem:(C,10):6090:6090,adj:95(1),15(1),112(1),179(1),
	// 179:20627,3 1-mem:(C,11):20627:20627,adj:178(1),151(1),225(1),226(1),
	// 180:6483,3 1-mem:(F,4):6483:6483,adj:94(1),67(1),197(1),195(1),
	// 182:8270,3 1-mem:(B,5):8270:8270,adj:72(1),102(1),245(1),194(1),
	// 183:2894,3 1-mem:(O,2):2894:2894,adj:93(1),279(1),273(1),184(1),
	// 184:844,3 1-mem:(O,1):844:844,adj:183(1),187(1),191(1),
	// 185:592,3 1-mem:(K,1):592:592,adj:88(1),2(1),186(1),
	// 187:781,3 1-mem:(N,1):781:781,adj:91(1),184(1),279(1),
	// 188:3216,3 1-mem:(T,2):3216:3216,adj:11(1),3(1),189(1),
	// 189:3018,3 1-mem:(S,2):3018:3018,adj:188(1),31(1),272(1),190(1),
	// 190:3086,3 1-mem:(R,2):3086:3086,adj:189(1),12(1),85(1),89(1),
	// 191:906,3 2-mem:(Q,1):906:906,(P,1),adj:89(1),85(1),184(1),273(1),
	// 192:10395,3 1-mem:(C,6):10395:10395,adj:245(1),246(1),194(1),193(1),
	// 194:8332,3 1-mem:(C,5):8332:8332,adj:192(1),182(1),44(1),196(1),
	// 195:6414,3 1-mem:(E,4):6414:6414,adj:180(1),181(1),177(1),198(1),
	// 196:8403,3 1-mem:(D,5):8403:8403,adj:193(1),194(1),181(1),198(1),
	// 197:8526,3 1-mem:(F,5):8526:8526,adj:121(1),180(1),208(1),198(1),
	// 198:8460,3 1-mem:(E,5):8460:8460,adj:132(1),197(1),196(1),195(1),
	// 199:19211,3 1-mem:(N,10):19211:19211,adj:14(1),157(1),155(1),230(1),
	// 200:17298,3 1-mem:(P,9):17298:17298,adj:107(1),35(1),206(1),203(1),
	// 202:15314,3 1-mem:(Q,8):15314:15314,adj:59(1),35(1),277(1),206(1),
	// 203:19340,3 1-mem:(P,10):19340:19340,adj:201(1),170(1),157(1),200(1),
	// 204:11155,3 1-mem:(P,6):11155:11155,adj:29(1),71(1),207(1),205(1),
	// 206:15248,3 1-mem:(P,8):15248:15248,adj:41(1),202(1),200(1),207(1),
	// 207:13195,3 1-mem:(P,7):13195:13195,adj:59(1),206(1),204(1),243(1),
	// 208:10574,3 1-mem:(F,6):10574:10574,adj:130(1),132(1),39(1),197(1),
	// 209:12499,3 1-mem:(D,7):12499:12499,adj:21(1),193(1),246(1),210(1),
	// 210:12562,3 1-mem:(E,7):12562:12562,adj:96(1),130(1),132(1),209(1),
	// 211:21203,3 1-mem:(M,11):21203:21203,adj:155(1),233(1),219(1),230(1),
	// 212:23066,3 1-mem:(J,12):23066:23066,adj:134(1),235(1),214(1),215(1),
	// 213:22859,3 1-mem:(F,12):22859:22859,adj:124(1),140(1),237(1),238(1),
	// 214:23115,3 1-mem:(K,12):23115:23115,adj:122(1),138(1),212(1),134(1),
	// 215:25099,3 1-mem:(J,13):25099:25099,adj:138(1),158(1),106(1),212(1),
	// 216:25296,3 1-mem:(M,13):25296:25296,adj:46(1),218(1),217(1),219(1),
	// 217:27341,3 1-mem:(M,14):27341:27341,adj:169(1),163(1),216(1),63(1),
	// 218:25227,3 1-mem:(L,13):25227:25227,adj:122(1),216(1),63(1),138(1),
	// 219:23245,3 1-mem:(M,12):23245:23245,adj:122(1),216(1),38(1),211(1),
	// 220:12747,3 1-mem:(H,7):12747:12747,adj:97(1),32(1),42(1),222(1),
	// 221:14862,3 1-mem:(J,8):14862:14862,adj:153(1),97(1),223(1),222(1),
	// 222:14795,3 1-mem:(H,8):14795:14795,adj:36(1),221(1),220(1),227(1),
	// 223:16915,3 1-mem:(J,9):16915:16915,adj:221(1),232(1),228(1),227(1),
	// 224:22740,3 1-mem:(D,12):22740:22740,adj:124(1),10(1),225(1),226(1),
	// 225:6090,3 1-mem:(C,12):6090:6090,adj:224(1),78(1),9(1),179(1),
	// 226:8138,3 1-mem:(D,11):8138:8138,adj:224(1),179(1),95(1),229(1),
	// 227:16846,3 1-mem:(H,9):16846:16846,adj:36(1),223(1),222(1),234(1),
	// 228:18962,3 1-mem:(J,10):18962:18962,adj:134(1),40(1),223(1),234(1),
	// 229:20750,3 1-mem:(E,11):20750:20750,adj:137(1),226(1),124(1),238(1),
	// 230:19150,3 1-mem:(M,10):19150:19150,adj:126(1),199(1),211(1),231(1),
	// 231:19091,3 1-mem:(L,10):19091:19091,adj:40(1),230(1),239(1),233(1),
	// 232:16974,3 1-mem:(K,9):16974:16974,adj:40(1),223(1),153(1),239(1),
	// 233:21131,3 1-mem:(L,11):21131:21131,adj:122(1),211(1),134(1),231(1),
	// 234:18891,3 1-mem:(H,10):18891:18891,adj:228(1),36(1),227(1),236(1),
	// 235:22990,3 1-mem:(H,12):22990:22990,adj:212(1),106(1),237(1),236(1),
	// 236:20942,3 1-mem:(H,11):20942:20942,adj:36(1),134(1),235(1),234(1),
	// 237:22926,3 1-mem:(G,12):22926:22926,adj:36(1),213(1),235(1),143(1),
	// 238:20818,3 1-mem:(F,11):20818:20818,adj:36(1),229(1),213(1),147(1),
	// 239:17035,3 1-mem:(L,9):17035:17035,adj:126(1),231(1),232(1),129(1),
	// 240:13069,3 1-mem:(N,7):13069:13069,adj:128(1),28(1),243(1),241(1),
	// 241:15116,3 1-mem:(N,8):15116:15116,adj:41(1),14(1),127(1),240(1),
	// 242:25420,3 1-mem:(O,13):25420:25420,adj:109(1),46(1),30(1),161(1),
	// 243:13134,3 1-mem:(O,7):13134:13134,adj:29(1),207(1),41(1),240(1),
	// 244:12371,3 1-mem:(B,7):12371:12371,adj:103(1),152(1),245(1),246(1),
	// 245:4042,3 1-mem:(B,6):4042:4042,adj:192(1),244(1),90(1),182(1),
	// 246:12428,3 1-mem:(C,7):12428:12428,adj:125(1),209(1),244(1),192(1),
	// 247:35918,3 1-mem:(S,18):35918:35918,adj:13(1),8(1),65(1),250(1),
	// 248:33235,3 1-mem:(H,17):33235:33235,adj:76(1),117(1),54(1),249(1),
	// 249:6090,3 1-mem:(G,17):6090:6090,adj:248(1),131(1),254(1),253(1),
	// 250:4042,3 1-mem:(R,18):4042:4042,adj:247(1),101(1),115(1),251(1),
	// 251:35790,3 1-mem:(Q,18):35790:35790,adj:60(1),250(1),173(1),252(1),
	// 252:33742,3 1-mem:(Q,17):33742:33742,adj:27(1),115(1),251(1),57(1),
	// 253:33099,3 1-mem:(F,17):33099:33099,adj:249(1),175(1),53(1),108(1),
	// 254:35214,3 1-mem:(G,18):35214:35214,adj:76(1),249(1),176(1),175(1),
	// 255:37458,3 1-mem:(K,19):37458:37458,adj:77(1),119(1),256(1),
	// 256:35403,3 1-mem:(K,18):35403:35403,adj:136(1),255(1),261(1),257(1),
	// 257:35340,3 1-mem:(J,18):35340:35340,adj:256(1),77(1),76(1),117(1),
	// 258:35534,3 1-mem:(M,18):35534:35534,adj:119(1),136(1),111(1),259(1),
	// 259:35595,3 1-mem:(N,18):35595:35595,adj:171(1),258(1),265(1),260(1),
	// 260:35662,3 1-mem:(O,18):35662:35662,adj:60(1),259(1),172(1),264(1),
	// 261:33358,3 1-mem:(K,17):33358:33358,adj:256(1),117(1),55(1),266(1),
	// 262:29459,3 1-mem:(N,15):29459:29459,adj:169(1),165(1),163(1),268(1),
	// 263:31571,3 1-mem:(O,16):31571:31571,adj:165(1),268(1),269(1),264(1),
	// 264:33619,3 1-mem:(O,17):33619:33619,adj:27(1),263(1),260(1),265(1),
	// 265:33548,3 1-mem:(N,17):33548:33548,adj:264(1),259(1),111(1),268(1),
	// 266:33420,3 1-mem:(L,17):33420:33420,adj:111(1),136(1),261(1),168(1),
	// 267:31442,3 1-mem:(M,16):31442:31442,adj:169(1),111(1),168(1),268(1),
	// 268:8138,3 1-mem:(N,16):8138:8138,adj:263(1),265(1),267(1),262(1),
	// 269:31630,3 1-mem:(P,16):31630:31630,adj:27(1),57(1),263(1),104(1),
	// 270:31756,3 1-mem:(R,16):31756:31756,adj:13(1),115(1),57(1),20(1),
	// 271:4302,3 1-mem:(D,3):4302:4302,adj:75(1),177(1),181(1),99(1),
	// 272:5202,3 1-mem:(S,3):5202:5202,adj:189(1),11(1),66(1),12(1),
	// 273:2955,3 1-mem:(P,2):2955:2955,adj:183(1),85(1),12(1),191(1),
	// 277:15372,3 1-mem:(R,8):15372:15372,adj:80(1),13(1),20(1),202(1),
	// 278:6924,3 1-mem:(N,4):6924:6924,adj:19(1),18(1),93(1),284(1),
	// 279:2832,3 1-mem:(N,2):2832:2832,adj:183(1),93(1),187(1),280(1),
	// 281:9298,3 1-mem:(S,5):9298:9298,adj:139(1),66(1),285(1),282(1),
	// 282:4042,3 1-mem:(S,6):4042:4042,adj:34(1),84(1),281(1),286(1),
	// 283:9037,3 1-mem:(O,5):9037:9037,adj:205(1),29(1),19(1),284(1),
	// 284:8976,3 1-mem:(N,5):8976:8976,adj:283(1),278(1),28(1),81(1),
	// 285:9234,3 1-mem:(R,5):9234:9234,adj:281(1),205(1),66(1),286(1),
	// 286:11276,3 1-mem:(R,6):11276:11276,adj:282(1),80(1),71(1),285(1),
	// Total 264 nodes, with 329 members
	// Move[0]: Loc: R4, Type: Black, Num: 1,
	// Move[1]: Loc: D4, Type: White, Num: 2,
	// Move[2]: Loc: F3, Type: Black, Num: 3,
	// Move[3]: Loc: D6, Type: White, Num: 4,
	// Move[4]: Loc: Q16, Type: Black, Num: 5,
	// Move[5]: Loc: D17, Type: White, Num: 6,
	// Move[6]: Loc: L3, Type: Black, Num: 7, CapdBy: 23
	// Move[7]: Loc: C14, Type: White, Num: 8,
	// Move[8]: Loc: Q10, Type: Black, Num: 9,
	// Move[9]: Loc: K16, Type: White, Num: 10,
	// Move[10]: Loc: M17, Type: Black, Num: 11,
	// Move[11]: Loc: P5, Type: White, Num: 12,
	// Move[12]: Loc: Q6, Type: Black, Num: 13,
	// Move[13]: Loc: L5, Type: White, Num: 14,
	// Move[14]: Loc: J4, Type: Black, Num: 15,
	// Move[15]: Loc: M3, Type: White, Num: 16,
	// Move[16]: Loc: M2, Type: Black, Num: 17,
	// Move[17]: Loc: L2, Type: White, Num: 18, CapdBy: 20
	// Move[18]: Loc: K2, Type: Black, Num: 19,
	// Move[19]: Loc: L4, Type: White, Num: 20,
	// Move[20]: Loc: L1, Type: Black, Num: 21, FirstCap: 17
	// Move[21]: Loc: K3, Type: White, Num: 22,
	// Move[22]: Loc: J3, Type: Black, Num: 23,
	// Move[23]: Loc: L2, Type: White, Num: 24, CapdBy: 30 FirstCap: 6
	// Move[24]: Loc: N3, Type: Black, Num: 25, Ko: L3
	// Move[25]: Loc: R3, Type: White, Num: 26,
	// Move[26]: Loc: Q4, Type: Black, Num: 27,
	// Move[27]: Loc: Q3, Type: White, Num: 28,
	// Move[28]: Loc: P4, Type: Black, Num: 29,
	// Move[29]: Loc: P3, Type: White, Num: 30,
	// Move[30]: Loc: L3, Type: Black, Num: 31, CapdBy: 33 FirstCap: 23
	// Move[31]: Loc: O4, Type: White, Num: 32, Ko: L2
	// Move[32]: Loc: O3, Type: Black, Num: 33,
	// Move[33]: Loc: L2, Type: White, Num: 34, FirstCap: 30
	// Move[34]: Loc: S4, Type: Black, Num: 35, Ko: L3
	// Move[35]: Loc: Q5, Type: White, Num: 36,
}
