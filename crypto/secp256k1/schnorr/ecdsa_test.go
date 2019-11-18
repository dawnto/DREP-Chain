// Copyright (c) 2015-2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package schnorr

import (
	"bytes"
	"encoding/hex"
	"github.com/drep-project/DREP-Chain/crypto/sha3"
	"math/rand"
	"testing"

	"github.com/drep-project/DREP-Chain/crypto/secp256k1"
)

type SchorrSigningTestVectorHex struct {
	msg   string
	nonce string
	priv  string
	sig   string
}

// schnorrSigningTestVectors were produced using the testing functions
// implemented in libsecp256k1.
// https://github.com/bitcoin/secp256k1/blob/258720851e24e23c1036b4802a185850e258a105/src/modules/schnorr/tests_impl.h
var schnorrSigningTestVectors = []SchorrSigningTestVectorHex{
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "D1C4C30F60582323A609B56B92270181EB05C3E5AB3E19AE1F768C65C6D09A29", "714D90C991E5D26CBF5771D8A84D087200AAA3197C3217A702ED8D69EA714CAB", "0A3E13BFD0B64C120AA25D27E3CD87678154A4461CE0AD471273927A6459F0C6" + "B9A36629C110ECEEEBBD52E7A5D491BB10AF59C3C73285B9427D1254F28DC460"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "86D9A69D76C1435EDC35347B50B4F944D30EDF8B5CB8E897E95F2C1F1B72D3C3", "60A30BC3BC7CDED4F13C9E3F20F69B8F7B4AB70E60825AE053FC88A2E7046C1F", "D60EFA079B194592A5200C60438A3617691FDE1B5FBCF788D0943A4BB69592F1" + "66D469F48267AA71DAF4BA996BA2BF3A99858C4BF854E2CDFC8AB7E6571D6A8C"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "499FE87D281A8EDEC40D7C29202CA93F9E612760C689543897255CC3B543F2E9", "0C23A2A854DD57AC4773533E84039BA165CA1F79BE8019BDF9EA3173741C67E9", "3C3483E5CDAAF894261071A948B1E21906CEF0293D10A3D20325EA84CC129B32" + "FF07618FAD7BE485A5A1C15DD6EE5485058D03514259714E724879AABCD70C5D"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "258E4C1130CACA9D1CF5DDB589551374DC06491BED02A72CBB8FC7D211E1ED20", "D4045E55A0FED015E1E90934C092C146680090F5538752152F35DAD6ECE70A45", "5EC936DC757F57473A84383F11511B78DB25ABEC5D0DEAE76ECFB30B7A006D9E" + "073D1EFC0C026F41BFBAF9E98FCA42CC6E1946123BD30B5BD27039FFF2D5FC48"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "8E4E9709D63F40167327325CDB6595AD806CABA520FC678993EA872151F6F520", "FD986300F2EB2F3F004D878B8E05C0E8D423DD4B7F112E1396406A0180F14956", "6025A47D8E92D32C417AF55414E6DB6FDEA98D64271C98C5C7FE4C4A6AF75727" + "4106736F0AC7C8782223BCFCDDB2D9396EF5169AE74E81EBE66F1EF3B5982A77"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "F6945B4B5A036BD85915F617054066E16B95392002CE7E7A4B1CD8C2190191A2", "80B060ACE5CB0D9971FB6A8E9A1342D88D144B56FE21A24B9183ED4DACADD30B", "7CF3258EEFE3B837916C21D7E13E0A9363FF6F82444D849D2607DB805B19BF68" + "1B78E1E9C509076E1361A0CE3CE46F4E155EC269D19EDCA9685788728EDD9269"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "1A25DEBF234C081CBFA24EF6A61F59ABE76EA58A91B5D749C31B5FBDFA7D80E9", "AF32FAC33D85C1BFE52C29CE47A9E62CF9E4B7EB66E94DACEDBB61AE0733F826", "31D285DA09369C500E4ADA47D868720852176B813AB25A998E7518855780FF08" + "2C04CE5E1BB33A951496475B594A3B42E740DAC3A0EBED07766CE79FA53DFB56"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "1C158A9FFB6CA28F2C2B4DE2680CCFD166D7FE24ED4E46125FDA392B1B547B24", "DEB515BB6074FBDD04EF7F17F055C82C2BE78C21639AED5F80663A0E89C4C9FE", "77383212DABB1F13B23EEE672C3D2D0CA84A58A148CAD023BC7D27F4C2A1BBD0" + "E3A2DB11E8DDC5D8DA1C55C14214FCA9C876EA0635DC260EDD47C8EE6C4B0F11"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "A11FBCD2EFB92CD714C088A5DAA8BEEE5645B7B70A253EB49F6A08E8E64B2F6F", "C57B3F9CB4E72C757BF477BCC5C6CBFEE8F9BFF35FCDE781F647648955F89D70", "180997C4BC9C6BDA81580330CDC4826B15A6A0B591023F84E5A5CCD8503D0E0D" + "172EA8B127F3BEAE03FBB6E3E20A8324A270905854B4EC4F6734C10A0A87DEE3"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "98BFCC52C846EB49CD728C566F58A62FF2D046D7B4A61CBDE8A2A3607AF5484D", "9A77F52ECFCF560FD9DE0B353CDE562FB1674B5ADD569832808C7015197F2600", "B9CD85998407F3E190BE5144BDF15ADDA8178C1F38C5168A2567C7895698DA40" + "E86BA74E0ACB2C7B9045826921265DAAF31FA05DAE8B83A6F3354B6CB493E52F"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "48DEFD73566CCAC57B1CB1FE7143767646B045A57D84C6C43B68958EDEED0439", "18F462C1C3EC3492613BA5E6C2D12C7A271948DB943A6A7081C8F9E58C308BB2", "EFC76E59819889BAEC540F5D89D28A4404B7EB98D34F941A7CA1D8C525256876" + "679BA0F14E5B2C1F69E1010D1B4FD0265AF34A583059805D5A3D13BBD4C8C83E"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "54415DD3EA817C1338D3AFDBBB142A25C2658B7A0E002BEA64E8DB88F5D62AC0", "56469DE2D40B2D7F0012042EFC7647D265D62254500B3D2864D9830D06C0536C", "DAD150AD11E5814E5DA4451A5D6EEE409779BD47D57C4AE7A788536E7E9426A3" + "DE5F61C9E54CE432BAB0FB8A4D6AF3D81A54062108DC2D23B397B7566AAB3A7A"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "7AC0AD497AD4F45285970132B2F15F48A61D2102EE60630A0345D20577070995", "2B72777C7784D50D538810A30AE3BF430A8A9C15E8658ACAED4D41207BA804D1", "827560F2F92A0905097E381BD0E962421A9E43E105585F9CFC3FB2D321369293" + "E16F37244D106FB0FFF2CF4A239F65C2A01D8C3FAAD187E39509AB0B16C7B72F"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "D7B06ED329566E59C551D71C52E9878DC49A2955159C30C5F373B29858A7CB89", "E76806E6BFDDC81B642A0A880313F8A552F448035B82CB2A91DD6175A497147C", "B772E3BFADC01A5088B9637E2E3D7D3269531091B79FC48ABCEE5A887BD3A11B" + "8BFA61100172D404D4938E097EEFFB0508320657BFE54699CD7490F5DC939333"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "23154117CD21BD4FCCB706F32D0945C7C835C632F02C6AC350B4E6BFDAA7F5BC", "96B0E0189F88AB41C486A65957CE6DEEEEBF2382F5393F9DAD1E8C65ADABEF25", "4D46A811696D6ADFCEDB303C4CA5474912502278D17BDB991C90C7EF54943B14" + "B2AF61F481E885C0B73033730AA78C3E59E7BB01B976D3ADA6EC5620FF9C279C"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "38B077D9A2448A2103CEA6F273C7F3720296A9A1814D0B71947D8AC83D00A86F", "836E462F675F2C9FCBF338F4F11494FF3A1BCA5221F6914B8219C5B189D77D40", "869A3BD110A8D32B6BFD8EF14B914C58276462EFF9CF8CE08279DF8BDAD24593" + "E34709964F528A1F0DFC112CAA2C3D26FE51E875CB6A495B57DCAB6876193511"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "FBD4EC7D96CB1F88BCBFEAF6FB7DDEE0898A3E4246C60FD688784B0F3C86444B", "8C5762E2019C91F24C80F1B4B7382912B573B3E92A8B9BDE3FCF4F902D84874C", "C2A63FE9E15D19D5ED6232A5641C225B7E0F06B4D4F0C53750E7BD889DAE0B36" + "50C13CFD562D104DEFD08C3CB208A4EFDF2EA405795E61D37162BF374E1F9CFF"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "76598BC70790E15BD717E9044BCDBD1B02EE09B9A7359CB77DFD1BC324C6AC54", "0CA72F2D1207B08162535A2C57282FB9B73C842818BB14378DB2F0C202F44BDA", "058A67F66CFFB0C7F14A5F8D98B068264E02DF7D372B0A1658308D6694F3E2D8" + "3B1014276872CD2ED9D19D3CE03AA0AE87F86FF1AC643D2179947552F363FA31"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "AAFEBB955B36B528034142464309C23C2F252F1A1BA97AB34D54FBD35899D896", "6A64A71C9BA0B2C5AD94F490CD76682B47C5EA7EE3B3B43208ADBB855DA9629E", "75541ADBACFE1483A27B2C261FFAD609BD38BDC9FE2C951B1B1A53C0929973C9" + "AC5B33E168B832CD2B65E1A164BBBB45939F26696367461A41EBA76984FD70A9"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "A19FDD9C1A161DC1B272E9484FD60B80374765D89C60928230BA51C671986262", "7B82469A8B8D1E9A12756869E255B70DB0404DB187565D88DFD78618077826F3", "6E5221F9AB9634644B3447C56504BE4B2C6DEC1024DC5E486F0D11E8ABFF7266" + "7096BD1F7FC33C5E1F7AAD7DF06C43FAF8B7A05F312ADFABF586D516B988E880"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "D463F98114D5EFBF206441C3BCB07B5D0358B9ABCD7212700D1E85A3B59B4B14", "7507E3E1CD0CAF7E951B8F46C7A52717203AA6922F5B3AA4B813AE05B3BED616", "8CEE361B3E0FB92B24CBC9E239E5D80DA558354E0108715286EFEFDEF55275B0" + "A68CAFB8E5C8499EE3895F24C007AA9A43E708B8A2959F43A8A9F8C14E4B2B8C"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "F15B442705E64A36688564BB3E5DC9ED8469AC73CE1794E549D31C34A8CF7367", "2C755FB91D25C48F0340519A54DF0026F7F6DC06B8D9DB3BC50BBF6A82344CCA", "45E262112967F578ABC595C6D12A887D60BCEB91FEA1F83EE6488F89159451ED" + "944891E90094B663718799A3729D0423FA95D7D327E480BE1556DE1F9D8C29DB"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "27261595EF702F48336F47DD2317461C06157C0868BD0E18A27BFA393551E5EE", "3DA9A1F356C519B7E4B7ACD05A648D4BEE8539ACF4B60D291E68D76B357D308B", "019C854E1D900122F706721816DB3895C6772B9EE254F39B326895D299910AF6" + "F58A6281C7DECCD1104094E1D84A5C182D549B55F834DB610F7CA0B5CD74A718"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "FF553D981379C5046D1927B50DB2B8D24A0274C09F7F4FBDDF100C59398DCC0E", "6470BE1A8EA8E9328B5F5F83B64677F4E346AEFCFCE71D9725E062998E0DCF42", "CE38C701700CBFA7442BB4E51C0C6F22FCCE1F39D4B22ADF5D5A910CD8E22CDD" + "55DC9688ED2F2ABFCEEB40C8FD217F274F9C369C8ADEF703AF7772C967C06F84"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "F38B04C155684AF2C717223DCC91EBC828B2D5B09EE807AE45C131693EC63637", "A8A1BEF67F7A3847DB2A57F3724467EDCC6454D29CA0E399F2A492898E69E497", "3F95D1916590BC74B6FDC4F3E412ED29E040358F17DFDD209C0F4AE04C9E94EA" + "D8281CCD545F57D80EA420CFDE5F81E350333BED9A39E116F72B17635AB57A92"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "C97757D41207BECB21EC0B4EB4AF7450F06124F43701C0260A46734F7C3D67B0", "D585B47F24BB9BC50A555A7FF2E706140404BCE393B2ADBBD927D3B4608294E9", "66BD8B3E23E11BE273AFA3257E9EF2D979A957E5C6B97A7966CFF11B4B581E67" + "EECA4295F554E20DB9F59793130595AAD73641C813AA15028190130592D21E95"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "2B58BCED7852B0F4198A0802CA0887DFA819BD73AE54E8EEE55FA4C76B28BE29", "638E732221DF5962F9ABD1FC69375058D20215F3ED226F6CCEED2C6F9D91B56B", "8AE5EFDF614C21FD637A3B26BB5655E4922F2598826A10DE2218B8BF4B7328F2" + "611C67385E260B9DD75C59623A003ECAD7F73811BE8BC7F015574ECFEBE3A9DD"},
	{"304502210088BE0644191B935DB1CD786B43FF27798006578D8C908906B49E89", "0371B3651E7ABC533801290E1C7E1D91D9AEB000ECF44FD0500773699A1C95C7", "6B7439B47606111ABBDE3429B6BD2938F0CEE87E7507A265655BAEFE25B48A16", "F3353583290AF1CC35AD4929633A0044EE4E7422E260AF341F597D240257C51D" + "F2813DA1E04C64A03FBBD57834897D6DEC8D441D1DECF80746502081355BAA3C"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "FFFFFFFFFF03000000FCFFFFFFFFFFFFFF030000000000007CE0FFFFFFFFFFFF", "E0FFFFFFFF1F000000E0FFFFFFFFFFFFFFFFFFFFFFFFFFFF1F00E0FFFFFFFFFF", "D1637C2FFDBAF642250F8B54FEE34A98CB7DE641BA1E9EFA9A20EA874A5FDCFA" + "F0A7B370F9458C0562C0AEF18BBBB84D50BAA6533E05F2C11DDEFC1E4BB1A6C5"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "0CF8FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFFFFFFF070000F8FF3F0000", "FFFFFFFFFF07000000000000C001000000E0FF3F00F0FFFFFFFFFFFFFF1F0000", "C0DB1A31A4E93F49E53C48CC1766CCC51B9A214400CC6A784090E7E0409D49BE" + "A0CAE6C3073B305E89A810F64558CC636F1F59049D00D4D7824A8BC1CB0FBEC4"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "E1FFFFFFFFFFFF00F07F0000000000000000E0FFFFFF3F0000FEFF0FF8FFFFFF", "000000000000000000000000F0FFFFFF1F3000E0FF7FF8FFFFFFFFFF0F000000", "AE12E3F8D79D596E023A250D9DDDB150ED35509F6F4B80488551CD46F1CA9D8A" + "78ED5DB488348DA811E1E9A192AF4E250F1EE3DC46C2A035E19D39D8FEAF9FEB"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "9FFFFFFFFFFFFFFFFFFFFFFFDF00FCC1FFFFFFFFFFFFFFFFFFFFFFFFFFFF0F00", "FFFFFFFFFFFFFFFFFFFF1F0000000000000000F8FFFF0F000000E0FFFFFFFFFF", "FE984A34D9D7B673DFCA1BAADC1F39C546BAA222E253E66726CA3045CCCC948B" + "A1DF6CBA13DD65FDBAAD8017A5C1C95331DEE2C07FFCDFE204ACD3E70028CA11"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "000000C03FFCFFFF03000000FFFF7FF8FFFF1F000000FFFFFFFF0F0000FC1F00", "FEFFFFFFFF1F00000000000000000000000000FFFF3F00E0070000FCFFFFFF03", "4BF5AD0E0C37DCE99BC56CD2D3AB0FDDF1321AE371E4F7D5E113D1D26793C780" + "DF0E1E0B5D69D151DD08639E25171A73C0CEA6F983F6B93ABCD5CD187A1853D8"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "FFFFFCFFFFDF010000000000F0FFFF3F000000FC1F007800FEFFFFFFFFFF0700", "00000000C0FFFFFFFFFFFFFFFFFFFFFF00FEFF03000000000080007E00000000", "3CCCCA01404CF09A1E906A00F61C1559FC3C8F1D29516D24D6BBD94A51664E0D" + "44B9CBA23BEC4D7DDA36DEA03F8D6077057C01ED9E33105629A9CC776DBE8985"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "00E0FFFFFFFFFFFFFFFF03FE3F000000000000000000FF1F0000000000000000", "FF3F00000000FEFFFFFF3F0000F8FFFFFF0F00C0FFFFFFFFFFFF0700000000C0", "2E4B97107F07E264BB3E7EBA6E019DE8B5FF11F2E2F42C92EFEB50DE38E9EC5D" + "2B3912BE0B0035FC5C8EC889312E4B47AAF460AD4A16752CE5266F395622D8BD"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "000000F8FFFFFFFFFF010000000000300000000000F8FF00000080FFFFFFFFFF", "FFFFFF07C0FF07000000980FE0FFFFFF0100000000000FE0FF0F0000FCFFFFFF", "52CC053205817B8727F8AEB8508345B77FC22F3C5CBC0A473D5AF2AD27F6C155" + "1DA91F27C36FED4DD90F61DD36457014F335C0773A0FBF423416294AE78E43FB"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "FFFFFFFFFFFFFFFFFFFFFCFFFFFFE3FFFFFFFFFFFF0F00E0FFFFFF0000F0FF1F", "FFFFFFFFFF3F0000C0FFFFFFFFFFFFFF0F0000F0FF7F1E00000080FFFFFF07FC", "95F2E4ACECBF9603303FDB868FD3F9DCB1D4DFD6C311FE00EA26E65B740C2B9E" + "BD6038044DBE6EF5791991EE1E968BF9D25402EAD2497C212EC526EF4C6EEF3A"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "0F0080FFFFFFFF0F00000000000000C0FFFFFFFFFFFFFF07000000C0FFFFFFFF", "FFFFFFFFFFFF01E0FFFFFFFFFFFFFFFF03F8FF7F00F0FFFF0700000000000000", "D163552127BF351D7918F66741435865FC04694090A4A9B3FDB2BAC13462A05C" + "8457E0464AEE6B27036C39066F0E1DF3DF952DF258FED6A42B9450F0FF27D9EB"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "003C803F00FFFFFFFFFFFFFFFF7F0000000000FE0700000000000000000000FF", "FF01000000C0FFFFFFFFFF3F00E00300003C000000E0FFFFFFFFFFFFFFFFFFFF", "00846EE6532D8B30C8912117FCBEF293EE79E212BA4507F33C8B90AF8DE35C01" + "8450F88A1B0F7EBCCBC155C30A68D51E60FFF9DC769B36E47C1AE34EDE9FA37B"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "FFFFFFFFFFFFFFFFFFFFFFC1FFFFFFFFFF010000000080FFFFFFFF1F00C0FFFF", "FFFFFFFFFFFFFFFF010000C0FF81FF81FF100000000000000000000F0000F00F", "0F4C75E1A7EBA72A2D0F014F477FCC765E6EC36350F92EF75F08EEB4C1C3195F" + "18368C0E7B4510534443B8279C741D3F2D1B01766B051B452C3D0698D8E8D08D"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "FFFFFF0100FCFFFF7F0000000000003F0000000000F0FFFFFFFF7FE0FFFFFFFF", "FFFFFFFFFF7FFFFFFFFFFF070000000000000000E0FFFFFFFFFF030000FCFFC3", "C05CFF7CE3C072431F05213A27ECD81288C83CBE48C39C43633B54B9AC0E6890" + "BC580EC26EFDEB276062CFE6B978F55430F757438183AA63B5F4B5226AAF675F"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "FFFF000000FFFFFFFFE1FFFFFFFFFFFFFF030000000000000000F0FFFF00FFFF", "FFFF0700000000C0FFFFFFFFFFFFFFFFFF1F00C0FFFFFFFF0700FFFF03000000", "8CE4A5F3CD68B81548A2F7D59CB610FAE94AC52FA07ED40A769F08552E5C5BE0" + "54C46CAE9374DC4A10F8B6D36A1A9D7C4B5AB09415DDF7D229193D9349D511D1"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "0000000000FFFFFFFFFFFF030000000000000000F81FE0FF3F00FEFF0F0000FE", "000000C0FFFF7F0000F8FFFFFFFFFF07000000F0FFFFFFFFFFFFFFFFFFFFFF0F", "A338C9A7F05E31D8D6870232B88E57B5EDD92A85423A9C256BE52B0DA92579E5" + "6467E2D50FAA93D2A56150ABC03D9DF69590E3AD72C4CC7D0C2FAC5AFEFE9085"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "0000000000000000E0FFFF0700F8FFFFFFFF1F00000000000000000000000000", "FFFF0F00E0FF1F00000000000000F0FFFFFFFFFFFF010000000080FFFFFFFFFF", "FFC45EC42040A6F58B402AEA3C18E40614842BA7EFA73FF843B896FB426CCFA7" + "AF9A01F69B5CC37977014A83BCB28C1E8847B897159B6069A6F988DE634CC2AB"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "07000000000000000000000000F8FF7F00FFFFFFFFFF00000000C0FFFFFFFFFF", "000000000000FC0F00000000000000000000000000FEFF0300C4FFFFFF3F80FF", "F25FD62A0081100E01A501FBF1EC42678C6B9BAE722ABD4F4979E99E8CC06E25" + "8FB87682E70ABC22E0EFE457C76D78BF293BCE128381AE5EDD03A4195CA8AD40"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "80FFF8FF3FFFFFFFFFFFFFFFFFFFFFFFFFFFF1FFFFFFFF1F0000FEFFFFFFFF3F", "F8FFFFFFFFFFFFFF0FFFFF3F00000000000000000000000000C0FF0700000000", "075A4EB0F330215AEECA9E9180F129BC6E2F0A53A49FFDE5A272C3F04DEDA081" + "8BEEBDC42915B1EE2106DB38D310E70B1CD1CCF228B0B3D3640B65289988DEDE"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "0100FEFF0100FFFFFFFF03000000FCFFFFFFFF3F0000E0FFFF01000000000000", "F8FF000000C0FFFF07000000000000000000FCFFFFFFFF03E001E0FFFFFFFFFF", "D1E8298DA8B147ABB128536F773BF456A3628719A237026BF5B458C8776F325B" + "85A448D3181B429C8AAC59343DD3E4DEBDBE01C44B8E6CE03C4F5FB4020DEBB9"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "00FEFFFFFF7F00FFFFFFFFFFFFFFFFFFFFFF0700FCFFFF00000000000080FF80", "0100000000000000000000000080FF1F0000C00300F0FFFF0700000000000000", "144B8536728044AEF88AF33D3087CCB1D1F50C41A5C972361469818A6FA498D4" + "9A536497E36544F704560801CE46DA59713B9B7287E100F77A69B9E88BE95B9F"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "00000000E0FFFFFFFF0100000000FFFFFF7F0000C0FFFFFF3F0080FFFFFFFFFF", "0000000000F0FFFF0FF8FF0300000000000000C0FFFFFFFF010000F8FFFF1F00", "ACAA7BE8ED867701F671DB71020483705480618CEF52BD96C0306ECED854702B" + "BF8CB65D50E21B0FAB21D6AD11DDAA0183CFE29F2FE348DF4A3E73A07111E313"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "000000000000000000C0FF0FFCFFFFFFFFFFFFFFFFFFFF1F0000000000000000", "0080FFFFFFFFFF000000C0FFFFFFFF1F00FCFFFFFFFFFF7F000000E0FFFF00E0", "F71708DB249A209F44771F98E3D659F5139823835B188982EE195A50D3D860B4" + "5226533153B14AA89D5F7958B9799BB893BFADC9528C871087AC795ADE224595"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "00000000E0FF7F0000000000000080FF1FF8FFFFC1FFFFFF3F00FFFF03000000", "01000000000000000080FFFFFFFFFFFF3F0000FEFFFFFFFFFF0F0000F8FF1F00", "54A25E83137863A94778D33F031BAC8DC10F39913A8B3FD792660F1034C70DBD" + "2B62E4B73DBFF01B97401ECD2F03DF81090FEC1A59FA161687A2BAAAEEE514DD"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "80FFFFFFFFFFFF0300F8FFFFFFFF1F00000000C0FF0300001C0080FF1F7F0000", "FF7FC0FFFFFFFFFFFFFFFFFF030000FCFF3F0000F0FFFFFFFFFFFF0F00060000", "75161FC21626DD84BD58778207FC7F58506E11204A95A603D5E8E6CA575CE2F1" + "F5984BFA2C970641936FACD9C9778AAE02EC89483AF879D0BC063FB516744E62"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "000000FCFFFF0100000000000000C0FFFFFFFFFFFFFFFFFFF7FFFFFF1F000000", "00FCFF3FC0FFFF0100000000000000004000000000FF7F00000000E0FFFF3F00", "971FA2BC49EBA95D1CCC10706D34C809ECB6B6521B4B028B3029E104404E7DE8" + "1833F25C71C6A100CDCACB2E154A1AA2B140671A79D001913F4AFDF830C255C2"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "FF030000E0FFFFFF3F80FFFF017E00FCFF0100FCFFFFFF010000F0FF07FCFFFF", "000000000080FFFFFFFF03FCFFFFFFFF0700000000E0FF0F0000000000F8FFF0", "7655EEC4926FEC2B7D0ABFB1A63333E9B9CB893A9366CC597090F55CD0D2602F" + "2827D561883C5D4924712D0EE305891E5869132F254F4AAD8E553E10A85975AE"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "000000FEFFFFFF7F0000008003000000000000000080FFFFFF07FEFFFFFF7F00", "FFFFFFFFFFFF3FF8FFFFFFFFFFFF0100000000FCFFFFFFFFFFFF0700E0FF0100", "EDE211DEE2487804585D84DAE822BF376C318CB44D6EF6EA9526BF1C26C6AB6A" + "5BF3BD71AA4A9D8990055BF9AA49AD82BF549B86D0E06B7D90F601E45C167F39"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "0700000000000000FCFFFFFFFFFFFF0084FFFF0000F8FFFFFF0700FEFFFFFFFF", "00001C00FFFF070000C0FF0100F0FF030000C0FFFFFFFF00FEFFFFFFFFFF0100", "2836BD0AD158AFA88ACAF13A0EC8DCE4414EEC0D282B2BA18820A712E252493B" + "E03F8FB67706EF558F2356B4385A25511613BB8D5132C91CD5084CAB5E56323B"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "0000000000FEFF0F0000000080FFFFFFFF0F000000FFFFFFFF0100F0FFFF0700", "000000000000000000FFFF03809FFFFFFFFFFFFF0F0000000000000000000000", "2A5837844FDEDC8C5FD40350737EB0412DFBEA43AF75736D2669D9E7A4341917" + "0EECFF456B3C3BA4AD1B74DC45C7FA0D77DDC52A60943687008C768839149E7C"},
	{"304402207C7BC9E2D115C4C5C3E50950E69B30A9810BD73946A6D23C4ACBFF2E", "FFFF01000000FEFFFFFFFFF3FFFFFFFFFFFF7F000000000000F07FF8FFFFFFFF", "FFFFFFFF80FFFF07FEFF7F00000000000000FC07FCFFFF0700000000E0FFFFFF", "93E083E71C14BA94479CBE92213A56FF2ECFF8F2B085B2B3AA5CC6E8FEFAEAC0" + "76444710354091BB4FE9A218E875885F81DD787241A766C4E0422C5C1D7AD271"},
}

type SchorrSigningTestVector struct {
	msg   []byte
	nonce []byte
	priv  []byte
	sig   []byte
}

func GetSigningTestVectors() []*SchorrSigningTestVector {
	var tvs []*SchorrSigningTestVector
	for _, v := range schnorrSigningTestVectors {
		msg, _ := hex.DecodeString(v.msg)
		nonce, _ := hex.DecodeString(v.nonce)
		priv, _ := hex.DecodeString(v.priv)
		sig, _ := hex.DecodeString(v.sig)
		lv := SchorrSigningTestVector{msg, nonce, priv, sig}
		tvs = append(tvs, &lv)
	}

	return tvs
}

// Horribly broken hash function. Do not use for anything but tests.
func testSchnorrHash(msg []byte) []byte {
	h32 := make([]byte, scalarSize)

	j := 32
	for i := 0; i < 32; i++ {
		h32[i] = msg[i] ^ msg[j]
		j++
	}

	return h32
}

func TestSchnorrSigning(t *testing.T) {
	tRand := rand.New(rand.NewSource(54321))
	tvs := GetSigningTestVectors()
	for _, tv := range tvs {
		_, pubkey := secp256k1.PrivKeyFromBytes(tv.priv)

		sig, err :=
			schnorrSign(tv.msg, tv.priv, tv.nonce, nil, nil,
				testSchnorrHash)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		cmp := bytes.Equal(sig.Serialize()[:], tv.sig[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		// Make sure they verify too while we're at it.
		_, err = schnorrVerify(sig.Serialize(), pubkey, tv.msg,
			testSchnorrHash)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		// See if we can recover the public keys OK.
		var pkRecover *secp256k1.PublicKey
		pkRecover, _, err = schnorrRecover(sig.Serialize(), tv.msg,
			testSchnorrHash)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		cmp = bytes.Equal(pubkey.Serialize()[:], pkRecover.Serialize()[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		// Screw up the signature at a random bit and make sure that breaks it.
		sigBad := sig.Serialize()
		pos := tRand.Intn(63)
		bitPos := tRand.Intn(7)
		sigBad[pos] ^= 1 << uint8(bitPos)
		_, err = schnorrVerify(sigBad, pubkey, tv.msg,
			testSchnorrHash)
		if err == nil {
			t.Fatalf("expected an error, got %v", err)
		}

		// Make sure it breaks pubkey recovery too.
		var valid bool
		pkRecover, valid, err = schnorrRecover(sigBad, tv.msg,
			testSchnorrHash)
		if valid {
			cmp = bytes.Equal(pubkey.Serialize()[:], pkRecover.Serialize()[:])
			if cmp {
				t.Fatalf("expected %v, got %v", false, cmp)
			}
		} else {
			if err == nil {
				t.Fatalf("expected an error, got %v", err)
			}
		}
	}
}

func randPrivKeyList(i int) []*secp256k1.PrivateKey {
	r := rand.New(rand.NewSource(54321))

	privKeyList := make([]*secp256k1.PrivateKey, i)
	for j := 0; j < i; j++ {
		for {
			bIn := new([32]byte)
			for k := 0; k < scalarSize; k++ {
				randByte := r.Intn(255)
				bIn[k] = uint8(randByte)
			}

			pks, _ := secp256k1.PrivKeyFromBytes(bIn[:])
			if pks == nil {
				continue
			}

			// No duplicates allowed.
			if j > 0 &&
				(bytes.Equal(pks.Serialize(), privKeyList[j-1].Serialize())) {
				r.Seed(int64(j) + r.Int63n(12345))
				continue
			}
			privKeyList[j] = pks
			r.Seed(int64(j) + 54321)
			break
		}
	}

	return privKeyList
}

type SignatureVerParams struct {
	pubkey *secp256k1.PublicKey
	msg    []byte
	sig    *Signature
}

func randSigList(i int) []*SignatureVerParams {
	r := rand.New(rand.NewSource(54321))

	privKeyList := make([]*secp256k1.PrivateKey, i)
	for j := 0; j < i; j++ {
		for {
			bIn := new([32]byte)
			for k := 0; k < scalarSize; k++ {
				randByte := r.Intn(255)
				bIn[k] = uint8(randByte)
			}

			pks, _ := secp256k1.PrivKeyFromBytes(bIn[:])
			if pks == nil {
				continue
			}
			privKeyList[j] = pks
			r.Seed(int64(j) + 54321)
			break
		}
	}

	msgList := make([][]byte, i)
	for j := 0; j < i; j++ {
		m := make([]byte, 32)
		for k := 0; k < scalarSize; k++ {
			randByte := r.Intn(255)
			m[k] = uint8(randByte)
		}
		msgList[j] = m
		r.Seed(int64(j) + 54321)
	}

	sigsList := make([]*Signature, i)
	for j := 0; j < i; j++ {
		r, s, err := Sign(privKeyList[j], msgList[j])
		if err != nil {
			panic("sign failure")
		}
		sig := &Signature{r, s}
		sigsList[j] = sig
	}

	sigStructList := make([]*SignatureVerParams, i)
	for j := 0; j < i; j++ {
		ss := new(SignatureVerParams)
		pkx, pky := privKeyList[j].Public()
		ss.pubkey = secp256k1.NewPublicKey(pkx, pky)
		ss.msg = msgList[j]
		ss.sig = sigsList[j]
		sigStructList[j] = ss
	}

	return sigStructList
}

// Use our actual hashing algorithm here.
func TestSignaturesAndRecovery(t *testing.T) {
	r := rand.New(rand.NewSource(54321))

	numSigs := 128
	sigList := randSigList(numSigs)

	for _, tv := range sigList {
		pubkey := tv.pubkey
		sig := tv.sig

		// Make sure we can verify the original signature.
		_, err := schnorrVerify(sig.Serialize(), pubkey, tv.msg, sha3.Hash256)
		if err != nil {
			t.Fatalf("expected an error, got %v", err)
		}

		ok := Verify(pubkey, tv.msg, sig.R, sig.S)
		if !ok {
			t.Fatalf("expected %v, got %v", true, ok)
		}

		// See if we can recover the public keys OK.
		var pkRecover *secp256k1.PublicKey
		pkRecover, _, err = schnorrRecover(sig.Serialize(), tv.msg, sha3.Hash256)
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}

		cmp := bytes.Equal(pubkey.Serialize()[:], pkRecover.Serialize()[:])
		if !cmp {
			t.Fatalf("expected %v, got %v", true, cmp)
		}

		// Screw up the signature at some random bits and make sure
		// that breaks it.
		numBadBits := r.Intn(2)
		sigBad := sig.Serialize()
		// (numBadBits*2)+1 --> always odd so at least one bit is different
		for i := 0; i < (numBadBits*2)+1; i++ {
			pos := r.Intn(63)
			bitPos := r.Intn(7)
			sigBad[pos] ^= 1 << uint8(bitPos)
		}
		_, err = schnorrVerify(sigBad, pubkey, tv.msg, sha3.Hash256)
		if err == nil {
			t.Fatalf("expected an error, got %v", err)
		}

		// Make sure it breaks pubkey recovery too.
		var valid bool
		pkRecover, valid, err = schnorrRecover(sigBad, tv.msg,
			testSchnorrHash)
		if valid {
			cmp := bytes.Equal(pubkey.Serialize()[:], pkRecover.Serialize()[:])
			if cmp {
				t.Fatalf("expected %v, got %v", false, cmp)
			}
		} else {
			if err == nil {
				t.Fatalf("expected an error, got %v", err)
			}
		}
	}
}

func benchmarkSigning(b *testing.B) {
	r := rand.New(rand.NewSource(54321))
	msg := []byte{
		0xbe, 0x13, 0xae, 0xf4,
		0xe8, 0xa2, 0x00, 0xb6,
		0x45, 0x81, 0xc4, 0xd1,
		0x0c, 0xf4, 0x1b, 0x5b,
		0xe1, 0xd1, 0x81, 0xa7,
		0xd3, 0xdc, 0x37, 0x55,
		0x58, 0xc1, 0xbd, 0xa2,
		0x98, 0x2b, 0xd9, 0xfb,
	}

	numKeys := 1024
	privKeyList := randPrivKeyList(numKeys)

	for n := 0; n < b.N; n++ {
		randIndex := r.Intn(numKeys - 1)
		_, _, err := Sign(privKeyList[randIndex], msg)
		if err != nil {
			panic("sign failure")
		}
	}
}

func BenchmarkSigning(b *testing.B) { benchmarkSigning(b) }

func benchmarkVerification(b *testing.B) {
	r := rand.New(rand.NewSource(54321))

	numSigs := 1024
	sigList := randSigList(numSigs)

	for n := 0; n < b.N; n++ {
		randIndex := r.Intn(numSigs - 1)
		ver := Verify(sigList[randIndex].pubkey,
			sigList[randIndex].msg,
			sigList[randIndex].sig.R,
			sigList[randIndex].sig.S)
		if !ver {
			panic("made invalid sig")
		}
	}
}

func BenchmarkVerification(b *testing.B) { benchmarkVerification(b) }
