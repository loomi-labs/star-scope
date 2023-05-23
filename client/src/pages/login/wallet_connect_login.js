export function wallet_connect_login(uri) {
    try {
        if (walletConnect) {
            if (!uri) {
                return {result: "", error: "No URI."}
            }

            if (walletConnect.isMobile()) {
                if (walletConnect.isAndroid()) {
                    // Save the mobile link.
                    walletConnect.saveMobileLinkInfo({
                        name: "Keplr",
                        href: "intent://wcV1#Intent;package=com.chainapsis.keplr;scheme=keplrwallet;end;",
                    });

                    return {result: `intent://wcV1?${uri}#Intent;package=com.chainapsis.keplr;scheme=keplrwallet;end;`, error: ""}
                } else {
                    // Save the mobile link.
                    walletConnect.saveMobileLinkInfo({
                        name: "Keplr",
                        href: "keplrwallet://wcV1",
                    });

                    return {result: `keplrwallet://wcV1?${uri}`, error: ""}
                }
            }
            return {result: "", error: "Is not mobile."}
        }
        return {result: "", error: "Error loading walletConnect."}
    } catch (error) {
        return {result: "", error: error.toString()}
    }
}
