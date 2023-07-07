import '../styles/globals.css';
import type {AppProps} from 'next/app';
import {ChainProvider, defaultTheme} from '@cosmos-kit/react';
import {ChakraProvider} from '@chakra-ui/react';
import {wallets as keplrWallets} from '@cosmos-kit/keplr';
import {wallets as cosmostationWallets} from '@cosmos-kit/cosmostation';
import {wallets as leapWallets} from '@cosmos-kit/leap';

import {assets, chains} from 'chain-registry';
import {getSigningCosmosClientOptions} from 'interchain';

import {SignerOptions} from '@cosmos-kit/core';
import {Chain} from '@chain-registry/types';

function CreateCosmosApp({Component, pageProps}: AppProps) {
  const signerOptions: SignerOptions = {
    // @ts-ignore
    signingStargate: (_chain: Chain) => {
      return getSigningCosmosClientOptions();
    },
  };

  return (
    <ChakraProvider theme={defaultTheme}>
      <ChainProvider
        chains={chains}
        assetLists={assets}
        // @ts-ignore
        // wallets={[...keplrWallets, ...cosmostationWallets, ...leapWallets]}
        wallets={[keplrWallets[0], cosmostationWallets[0], ...leapWallets]}
        walletConnectOptions={{
          signClient: {
            projectId: 'a8510432ebb71e6948cfd6cde54b70f7',
            relayUrl: 'wss://relay.walletconnect.org',
            metadata: {
              name: 'Starscope',
              description: 'Starscope Network',
              url: 'https://starscope.network',
              icons: [],
            },
          },
        }}
        wrappedWithChakra={false}
        signerOptions={signerOptions}
      >
        <Component {...pageProps}/>
      </ChainProvider>
    </ChakraProvider>
  );
}

export default CreateCosmosApp;
