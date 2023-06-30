import {useState} from 'react';
import {useChain} from '@cosmos-kit/react';
import {AminoSignResponse, StdSignDoc} from '@cosmjs/amino';
import BigNumber from 'bignumber.js';

import {Container, useColorMode,} from '@chakra-ui/react';
import {chainassets, chainName, coin,} from '../config';
import {WalletSection,} from '../components';

import {cosmos} from 'interchain';
import {OfflineSigner} from "@cosmjs/proto-signing";
import {SignOptions} from "@cosmos-kit/core/types/types/wallet";

const library = {
  title: 'Interchain',
  text: 'Interchain',
  href: 'https://github.com/cosmology-tech/interchain',
};

const signMsg = (
  getOfflineSigner: () => Promise<OfflineSigner>,
  signAmino: (signer: string, signDoc: StdSignDoc, signOptions?: SignOptions) => Promise<AminoSignResponse>,
  setResp: (resp: string) => any,
  address: string
) => {
  return async () => {
    const offlineSigner = await getOfflineSigner();
    if (!offlineSigner || !address) {
      console.error('stargateClient undefined or address undefined.');
      return;
    }

    const signMsg = {
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
            signer: address,
            data: btoa("Hello".toLowerCase()),
          },
        },
      ],
      memo: "",
    }
    const response = await signAmino(address, signMsg,
      {preferNoSetFee: true, preferNoSetMemo: true, disableBalanceCheck: true});
    setResp(JSON.stringify(response, null, 2));
  };
};

export default function Home() {
  const {colorMode, toggleColorMode} = useColorMode();

  const {getOfflineSigner, signAmino, getSigningStargateClient, address, status, getRpcEndpoint} = useChain(chainName);

  const [balance, setBalance] = useState(new BigNumber(0));
  const [isFetchingBalance, setFetchingBalance] = useState(false);
  const [resp, setResp] = useState('');
  const getBalance = async () => {
    if (!address) {
      setBalance(new BigNumber(0));
      setFetchingBalance(false);
      return;
    }

    let rpcEndpoint = await getRpcEndpoint();

    if (!rpcEndpoint) {
      console.log('no rpc endpoint — using a fallback');
      rpcEndpoint = `https://rpc.cosmos.directory/${chainName}`;
    }

    // get RPC client
    const client = await cosmos.ClientFactory.createRPCQueryClient({
      rpcEndpoint,
    });

    // fetch balance
    const balance = await client.cosmos.bank.v1beta1.balance({
      address,
      denom: chainassets?.assets[0].base as string,
    });

    // Get the display exponent
    // we can get the exponent from chain registry asset denom_units
    const exp = coin.denom_units.find((unit) => unit.denom === coin.display)
      ?.exponent as number;

    // show balance in display values by exponentiating it
    const a = new BigNumber(balance.balance.amount);
    const amount = a.multipliedBy(10 ** -exp);
    setBalance(amount);
    setFetchingBalance(false);
  };

  return (
    <Container maxW="5xl" py={10}>
      <WalletSection
        handleSingMsg={signMsg(
          getOfflineSigner as () => Promise<OfflineSigner>,
          signAmino as () => Promise<AminoSignResponse>,
          setResp as () => any,
          address as string
        )}
      />
    </Container>
  );
}
