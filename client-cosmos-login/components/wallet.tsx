import { useChain } from '@cosmos-kit/react';
import {
  Box,
  Center,
  Grid,
  GridItem,
  Icon,
  Stack,
  useColorModeValue,
  Text,
} from '@chakra-ui/react';
import {MouseEventHandler, useEffect, useState} from 'react';
import { FiAlertTriangle } from 'react-icons/fi';
import {
  Astronaut,
  Error,
  Connected,
  ConnectedShowAddress,
  ConnectedUserInfo,
  Connecting,
  ConnectStatusWarn,
  CopyAddressBtn,
  Disconnected,
  NotExist,
  Rejected,
  RejectedWarn,
  WalletConnectComponent,
  ChainCard,
} from '../components';
import { chainName } from '../config';
import {WalletStatus} from "@cosmos-kit/core";

export const WalletSection = ({handleSingMsg}: {handleSingMsg: () => void}) => {
  const {
    connect,
    openView,
    status,
    username,
    address,
    message,
    wallet,
    chain: chainInfo,
    logoUrl,
  } = useChain(chainName);

  const chain = {
    chainName,
    label: chainInfo.pretty_name,
    value: chainName,
    icon: logoUrl,
  };

  const [userInteracted, setUserInteracted] = useState(false);

  // Events
  const onClickConnect: MouseEventHandler = async (e) => {
    e.preventDefault();
    await connect();
    setUserInteracted(true);
  };

  const onClickOpenView: MouseEventHandler = (e) => {
    e.preventDefault();
    openView();
    setUserInteracted(true);
    handleSingMsg();
  };

  useEffect(() => {
    if (status === WalletStatus.Connected && userInteracted) {
      handleSingMsg();
    }
  }, [status, userInteracted]);

  // Components
  const connectWalletButton = (
    <WalletConnectComponent
      walletStatus={status}
      disconnect={
        <Disconnected buttonText="Login with wallet" onClick={onClickConnect} />
      }
      connecting={<Connecting />}
      connected={
        <Connected buttonText={'Login with wallet'} onClick={onClickOpenView} />
      }
      rejected={<Rejected buttonText="Reconnect" onClick={onClickConnect} />}
      error={<Error buttonText="Change Wallet" onClick={onClickOpenView} />}
      notExist={
        <NotExist buttonText="Install Wallet" onClick={onClickOpenView} />
      }
    />
  );

  const connectWalletWarn = (
    <ConnectStatusWarn
      walletStatus={status}
      rejected={
        <RejectedWarn
          icon={<Icon as={FiAlertTriangle} mt={1} />}
          wordOfWarning={`${wallet?.prettyName}: ${message}`}
        />
      }
      error={
        <RejectedWarn
          icon={<Icon as={FiAlertTriangle} mt={1} />}
          wordOfWarning={`${wallet?.prettyName}: ${message}`}
        />
      }
    />
  );

  return (
    <Center>
      <Grid
        w="full"
        maxW="sm"
        templateColumns="1fr"
        rowGap={4}
        alignItems="center"
        justifyContent="center"
      >
        <GridItem px={6}>
          <Stack
            justifyContent="center"
            alignItems="center"
          >
            <Box w="full" maxW={{ base: 52, md: 64 }}>
              {connectWalletButton}
            </Box>
            {connectWalletWarn && <GridItem>{connectWalletWarn}</GridItem>}
          </Stack>
        </GridItem>
      </Grid>
    </Center>
  );
};
