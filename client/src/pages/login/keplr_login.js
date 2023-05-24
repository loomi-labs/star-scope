export async function keplr_login() {
    try {
        if (window.keplr) {
            const CHAIN_ID = "cosmoshub-4";
            await window.keplr.enable(CHAIN_ID);

            const offlineSigner = await window.getOfflineSigner(CHAIN_ID);
            const keplrAccounts = await offlineSigner.getAccounts();

            if (keplrAccounts[0].address !== "") {
                const signedMessage = await window.keplr.signAmino(CHAIN_ID, keplrAccounts[0].address, {
                    chain_id: "",
                    account_number: "0",
                    sequence: "0",
                    fee: {
                        gas: "0",
                        amount: [],
                    },
                    msgs: [
                        {
                            type: "sign/MsgSignData",
                            value: {
                                signer: keplrAccounts[0].address,
                                data: btoa("Hello".toLowerCase()),
                            },
                        },
                    ],
                    memo: "",
                })
                return {result: JSON.stringify(signedMessage), error: ""}
            }
        } else {
            return {result: "", error: "Keplr extension is not installed."}
        }
    } catch (error) {
        return {result: "", error: error.toString()}
    }
}
