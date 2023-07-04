import {
    isMobile as isMobileWC,
    isAndroid as isAndroidWC,
    saveMobileLinkInfo as saveMobileLinkInfoWC,
} from "@walletconnect/browser-utils";

export function isMobile(uri) {
    return isMobileWC();
}

export function isAndroid(uri) {
    return isAndroidWC();
}

export function saveMobileLinkInfo(uri) {
    return saveMobileLinkInfoWC(uri);
}