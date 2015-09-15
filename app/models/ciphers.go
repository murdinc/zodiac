package models

//import "github.com/revel/revel"

// Solved Cipher
const Cipher408 = "HER>plvVPk|1LTG3dNp+B"

// Unsolved Cipher
const Cipher340 = "HER>plvVPk|1LTG3dNp+B7$O%DWY.<^Kf6ByIcM+UZGW76L#$HJSpp#vl!^V4pO++RK3&@M+9tjd|0FP+P2k/p!RvFlO-^dCkF>2D7$0+Kq%i3UcXGV.9L|7G3Jfj$O+&NY9+*L@d<M+b+ZR3FBcyA52K-9lUV+vJ+Op#<FBy-U+R/0tE|DYBpbTMKO2<clRJ|^0T0M.+PBF95@Sy$+N|0FBc7i!RlGFNvf030b.cV0t++yBX1^I2@CE>VUZ0-+|c.49BK7Opv.fMqG3RcT+L03C<+FlWB|6L++6WC9WcPOSHT/76p|FkdW<#tB&YOB^-Cc>MDHNpkS9ZO!A|Ki+"

//const Cipher340Array

type Cipher struct {
    PerLine int         // Number of symbols per line
    Cipher  string      // The actual cipher
}

func (c *Cipher) GetCipher() {

}
