import {useState} from 'react';
import {useChain} from '@cosmos-kit/react';
import {AminoSignResponse, StdSignDoc} from '@cosmjs/amino';

import {Container,} from '@chakra-ui/react';
import {chainName,} from '../config';
import {WalletSection,} from '../components';
import {SignOptions} from "@cosmos-kit/core/types/types/wallet";

const signMsg = (
  signAmino: (signer: string, signDoc: StdSignDoc, signOptions?: SignOptions) => Promise<AminoSignResponse>,
  setResp: (resp: string) => any,
  address: string
) => {
  return async () => {
    if (!address) {
      console.error('address undefined.');
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
            data: btoa("hello"),
          },
        },
      ],
      memo: "",
    }
    try {
      const response = await signAmino(address, signMsg,
        {preferNoSetFee: true, preferNoSetMemo: true, disableBalanceCheck: true});
      const result = JSON.stringify(response, null, 2)
      window.parent.postMessage(result, '*');
      setResp(result);
    } catch (e) {
      console.error(e);
      window.parent.postMessage(e, '*');
    }
  };
};

export default function Home() {

  const {signAmino, address, status} = useChain(chainName);

  const [resp, setResp] = useState('');

  return (
    <Container maxW="5xl" bg="#342335">
      <WalletSection
        handleSingMsg={signMsg(
          signAmino as () => Promise<AminoSignResponse>,
          setResp as () => any,
          address as string
        )}
      />
    </Container>
  );
}
