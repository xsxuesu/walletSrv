package hdaddr

import (
	"fmt"
	"github.com/foxnut/go-hdwallet"
)

var (
	mnemonic = "range sheriff try enroll deer over ten level bring display stamp recycle"
)


func GetBtcHdAddress(index uint32 )(string,string,error){
	master, err := hdwallet.NewKey(
		hdwallet.Mnemonic(mnemonic),
	)
	if err != nil {
		return "","",err
	}

	wallet, err := master.GetWallet(hdwallet.CoinType(hdwallet.BTC), hdwallet.AddressIndex(index))
	if err != nil {
		return "","",err
	}

	address, err := wallet.GetAddress()
	if err != nil {
		return "","",err
	}

	//addressP2WPKH, _ := wallet.GetKey().AddressP2WPKH()
	//addressP2WPKHInP2SH, _ := wallet.GetKey().AddressP2WPKHInP2SH()
	btcPrik,err  := wallet.GetKey().PrivateWIF(true)
	if err != nil {
		return "","",err
	}

	return address,btcPrik,nil
}


func GetEthHdAddress(index uint32 )(string,string,error){
	master, err := hdwallet.NewKey(
		hdwallet.Mnemonic(mnemonic),
	)
	if err != nil {
		return "","",err
	}

	wallet, err := master.GetWallet(hdwallet.CoinType(hdwallet.ETH),hdwallet.AddressIndex(index))
	if err != nil {
		return "","",err
	}
	address, err := wallet.GetAddress()
	if err != nil {
		return "","",err
	}
	prik := wallet.GetKey().PrivateHex()

	return address,prik,err
}

func GetUSDTHdAddress(index uint32 )(string,string,error){
	master, err := hdwallet.NewKey(
		hdwallet.Mnemonic(mnemonic),
	)
	if err != nil {
		return "","",err
	}

	wallet, err := master.GetWallet(hdwallet.CoinType(hdwallet.USDT),hdwallet.AddressIndex(index))
	if err != nil {
		return "","",err
	}
	address, err := wallet.GetAddress()
	if err != nil {
		return "","",err
	}
	prik ,err := wallet.GetKey().PrivateWIF(true)

	if err != nil {
		return "","",err
	}

	return address,prik,err
}


func CreateHdAddress() {
	master, err := hdwallet.NewKey(
		hdwallet.Mnemonic(mnemonic),
	)
	if err != nil {
		panic(err)
	}

	// BTC: 1
	// 14ULoqLmWb6jEFGpD2FG2nfqmKqLMXeiTo
	// 14ULoqLmWb6jEFGpD2FG2nfqmKqLMXeiTo
	// Kx111Zh68Zb9kYRytVVFP6VDC4WgeHZGjfuKPCc4LfUW48jspkf3
	// Kx111Zh68Zb9kYRytVVFP6VDC4WgeHZGjfuKPCc4LfUW48jspkf3
	// 10000
	// 14SJdsw1T3ijqU39mty6EAxR1L3wwoGjUz
	// 14SJdsw1T3ijqU39mty6EAxR1L3wwoGjUz
	// L4kHAwBqK8T5MWuUfgsRzZ3tPQdPg8siMQna5vLZbQYUydYeWZux
	// L4kHAwBqK8T5MWuUfgsRzZ3tPQdPg8siMQna5vLZbQYUydYeWZux
	wallet, _ := master.GetWallet(hdwallet.CoinType(hdwallet.BTC), hdwallet.AddressIndex(10000))
	address, _ := wallet.GetAddress()
	addressP2WPKH, _ := wallet.GetKey().AddressP2WPKH()
	addressP2WPKHInP2SH, _ := wallet.GetKey().AddressP2WPKHInP2SH()
	btcPrik,err  := wallet.GetKey().PrivateWIF(true)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("BTC: ", address, addressP2WPKH, addressP2WPKHInP2SH)
	fmt.Println("BTC prikey: ", btcPrik)



	// 0x37039021cBA199663cBCb8e86bB63576991A28C1
	// ETH: 0x37039021cBA199663cBCb8e86bB63576991A28C1
	// 0xd14faa1f995260a7479de7beb126df2d6cef857b3a351bc9b3e7fbdb925c9a2e
	// 0xd14faa1f995260a7479de7beb126df2d6cef857b3a351bc9b3e7fbdb925c9a2e
	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.ETH),hdwallet.AddressIndex(10000))
	address, _ = wallet.GetAddress()
	prik := wallet.GetKey().PrivateHex()
	fmt.Println("ETH: ", address)
	fmt.Println("ETH privatekey: ", prik)

	// 1000
	// 0x2BA9B41a7Fb983FfB3B8BEA90D04909DD484dBB6
	// 0x2BA9B41a7Fb983FfB3B8BEA90D04909DD484dBB6
	// 0x2d7fb66933f68595903f1036dd3d48d657fde16906c951177b121ce3dc5df5ad
	// 0x2d7fb66933f68595903f1036dd3d48d657fde16906c951177b121ce3dc5df5ad
	// 10000
	// 0x289B39d70EccC02001658282911F359Fb3662249
	// 0x289B39d70EccC02001658282911F359Fb3662249
	// 0xbc8b296042de62eed2694bf5fb4961ac0c26ad4f5d0a27d2cf6c56603c16d1eb
	// 0xbc8b296042de62eed2694bf5fb4961ac0c26ad4f5d0a27d2cf6c56603c16d1eb



	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.USDT),hdwallet.AddressIndex(0))

	address, _ = wallet.GetAddress()
	usdtPrik,err  := wallet.GetKey().PrivateWIF(true)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("USDT: ", address)
	fmt.Println("USDT prikey: ", usdtPrik)
}

