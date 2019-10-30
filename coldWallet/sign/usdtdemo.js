/**
 * usdt 离线签名
 * @param privateKey
 * @param changeAddress
 * @param changeAmount
 * @param toAddress
 * @param outputs
 * @param amount
 * @return
 */
public String sign(String privateKey, String changeAddress,Long changeAmount, String toAddress, List<Utxo> outputs,Long amount)
{
    MainNetParams
    network = MainNetParams.get();
    Transaction
    tran = new Transaction(MainNetParams.get());

    //这是比特币的限制最小转账金额，所以很多usdt转账会收到一笔0.00000546的btc
    tran.addOutput(Coin.valueOf(546L),Address.fromBase58(network, toAddress));

    //构建usdt的输出脚本 注意这里的金额是要乘10的8次方
    String
    usdtHex = "6a146f6d6e69" + String.format("%016x", 31) + String.format("%016x", amount);
    tran.addOutput(Coin.valueOf(0L),new Script(Utils.HEX.decode(usdtHex)));

    //如果有找零就添加找零
    if (changeAmount.compareTo(0L) >0)
    {
        tran.addOutput(Coin.valueOf(changeAmount), Address.fromBase58(network, changeAddress));
    }

    //先添加未签名的输入，也就是utxo
    for (Utxo output : outputs)
    {
        tran.addInput(Sha256Hash.wrap(output.getTxHash()), output.getVout(), new Script(HexUtil.decodeHex(output.getScriptPubKey()))).setSequenceNumber(TransactionInput.NO_SEQUENCE - 2);
    }

    //下面就是签名
    for (int i = 0;i < outputs.size();i++)
    {
        Utxo
        output = outputs.get(i);
        ECKey
        ecKey = DumpedPrivateKey.fromBase58(network, privateKey).getKey();
        TransactionInput
        transactionInput = tran.getInput(i);
        Script
        scriptPubKey = ScriptBuilder.createOutputScript(Address.fromBase58(network, output.getAddress()));
        Sha256Hash
        hash = tran.hashForSignature(i, scriptPubKey, Transaction.SigHash.ALL, false);
        ECKey.ECDSASignature
        ecSig = ecKey.sign(hash);
        TransactionSignature
        txSig = new TransactionSignature(ecSig, Transaction.SigHash.ALL, false);
        transactionInput.setScriptSig(ScriptBuilder.createInputScript(txSig, ecKey));
    }
    //这是签名之后的原始交易，直接去广播就行了
    String
    signedHex = HexUtil.encodeHexStr(tran.bitcoinSerialize());
    //这是交易的hash
    String
    txHash = HexUtil.encodeHexStr(Utils.reverseBytes(Sha256Hash.hash(Sha256Hash.hash(tran.bitcoinSerialize()))));
    return signedHex;
}