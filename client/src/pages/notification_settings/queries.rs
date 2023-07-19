use crate::components::messages::{
    create_error_msg_from_status, create_message, create_timed_message,
};
use crate::types::protobuf::grpc_settings::{
    Chain, RemoveChainRequest, RemoveWalletRequest, UpdateChainRequest, UpdateWalletRequest,
    ValidateWalletRequest, Wallet,
};
use crate::{InfoLevel, Services};
use sycamore::prelude::*;

pub async fn query_wallets(cx: Scope<'_>) -> Result<Vec<Wallet>, ()> {
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(());
    let response = services
        .grpc_client
        .get_settings_service()
        .get_wallets(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(result) = response {
        Ok(result.wallets)
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
        Err(())
    }
}
pub enum WalletUpdate {
    Funding,
    Staking,
    GovVotingReminder,
}

pub async fn update_existing_wallet(
    cx: Scope<'_>,
    wallet_sig: &Signal<Wallet>,
    update: WalletUpdate,
) {
    let wallet = create_ref(cx, wallet_sig.get_untracked());
    let notify_funding = if let WalletUpdate::Funding = update {
        !wallet.notify_funding
    } else {
        wallet.notify_funding
    };
    let notify_staking = if let WalletUpdate::Staking = update {
        !wallet.notify_staking
    } else {
        wallet.notify_staking
    };
    let notify_gov_voting_reminder = if let WalletUpdate::GovVotingReminder = update {
        !wallet.notify_gov_voting_reminder
    } else {
        wallet.notify_gov_voting_reminder
    };
    let request = UpdateWalletRequest {
        wallet_address: wallet.address.clone(),
        notify_funding,
        notify_staking,
        notify_gov_voting_reminder,
    };

    let result = update_wallet(cx, wallet_sig, request).await;
    if result.is_ok() {
        let msg = match update {
            WalletUpdate::Funding => {
                if notify_funding {
                    "You will be notified about funding events"
                } else {
                    "You will no longer be notified about funding events"
                }
            }
            WalletUpdate::Staking => {
                if notify_staking {
                    "You will be notified about staking events"
                } else {
                    "You will no longer be notified about staking events"
                }
            }
            WalletUpdate::GovVotingReminder => {
                if notify_gov_voting_reminder {
                    "You will be notified to vote on governance proposals"
                } else {
                    "You will no longer be notified to vote on governance proposals"
                }
            }
        };
        create_timed_message(cx, "Wallet updated", msg, InfoLevel::Success, 5);
    }
}

pub async fn update_wallet(
    cx: Scope<'_>,
    wallet: &Signal<Wallet>,
    update: UpdateWalletRequest,
) -> Result<(), ()> {
    let update_ref = create_ref(cx, update.clone());
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(update);
    let response = services
        .grpc_client
        .get_settings_service()
        .update_wallet(request)
        .await
        .map(|res| res.into_inner());

    if response.is_ok() {
        wallet.modify().notify_funding = update_ref.notify_funding;
        wallet.modify().notify_staking = update_ref.notify_staking;
        wallet.modify().notify_gov_voting_reminder = update_ref.notify_gov_voting_reminder;
        Ok(())
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
        Err(())
    }
}

#[derive(Debug, Clone, PartialEq)]
pub enum WalletValidation {
    Valid(Wallet),
    Invalid(String),
}

pub async fn query_validate_wallet(cx: Scope<'_>, address: String) -> WalletValidation {
    let services = use_context::<Services>(cx);
    let request = services
        .grpc_client
        .create_request(ValidateWalletRequest { address });
    let response = services
        .grpc_client
        .get_settings_service()
        .validate_wallet(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(response) = response {
        if response.is_valid {
            if response.is_supported {
                if response.is_already_added {
                    return WalletValidation::Invalid("Wallet already added".to_string());
                }
                if let Some(wallet) = response.wallet {
                    return WalletValidation::Valid(wallet);
                }
                create_message(cx, "Error", "Wallet not found", InfoLevel::Error);
                WalletValidation::Invalid("Wallet not found".to_string())
            } else {
                WalletValidation::Invalid("Chain is currently not supported".to_string())
            }
        } else {
            WalletValidation::Invalid("Invalid wallet address".to_string())
        }
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
        WalletValidation::Invalid("Error".to_string())
    }
}

pub async fn delete_wallet(cx: Scope<'_>, wallet_address: String) -> Result<(), ()> {
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(RemoveWalletRequest {
        wallet_address: wallet_address.clone(),
    });
    let response = services
        .grpc_client
        .get_settings_service()
        .remove_wallet(request)
        .await
        .map(|res| res.into_inner());

    if response.is_ok() {
        create_timed_message(
            cx,
            "Wallet deleted",
            format!("Wallet {} was deleted", wallet_address),
            InfoLevel::Success,
            5,
        );
        Ok(())
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
        Err(())
    }
}

pub async fn query_chains(cx: Scope<'_>) -> Result<Vec<Chain>, ()> {
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(());
    let response = services
        .grpc_client
        .get_settings_service()
        .get_chains(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(result) = response {
        Ok(result.chains)
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
        Err(())
    }
}

pub enum ChainUpdate {
    NewProposal,
    ProposalFinished,
}

pub async fn update_existing_chain(cx: Scope<'_>, chain_sig: &Signal<Chain>, update: ChainUpdate) {
    let chain = create_ref(cx, chain_sig.get_untracked());
    let notify_new_proposals = if let ChainUpdate::NewProposal = update {
        !chain.notify_new_proposals
    } else {
        chain.notify_new_proposals
    };
    let notify_proposal_finished = if let ChainUpdate::ProposalFinished = update {
        !chain.notify_proposal_finished
    } else {
        chain.notify_proposal_finished
    };
    let request = UpdateChainRequest {
        chain_id: chain.id,
        notify_new_proposals,
        notify_proposal_finished,
    };

    let result = update_chain(cx, chain_sig, request).await;
    if result.is_ok() {
        let msg = match update {
            ChainUpdate::NewProposal => {
                if notify_new_proposals {
                    "You will be notified about new governance proposals"
                } else {
                    "You will no longer be notified about new governance proposals"
                }
            }
            ChainUpdate::ProposalFinished => {
                if notify_proposal_finished {
                    "You will be notified about finished governance proposals"
                } else {
                    "You will no longer be notified about finished governance proposals"
                }
            }
        };
        create_timed_message(cx, "Chain updated", msg, InfoLevel::Success, 5);
    }
}

pub async fn update_chain(
    cx: Scope<'_>,
    chain: &Signal<Chain>,
    update: UpdateChainRequest,
) -> Result<(), ()> {
    let update_ref = create_ref(cx, update.clone());
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(update);
    let response = services
        .grpc_client
        .get_settings_service()
        .update_chain(request)
        .await
        .map(|res| res.into_inner());

    if response.is_ok() {
        chain.modify().notify_new_proposals = update_ref.notify_new_proposals;
        chain.modify().notify_proposal_finished = update_ref.notify_proposal_finished;
        Ok(())
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
        Err(())
    }
}

pub async fn delete_chain(cx: Scope<'_>, chain: Chain) -> Result<(), ()> {
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(RemoveChainRequest {
        chain_id: chain.id,
    });
    let response = services
        .grpc_client
        .get_settings_service()
        .remove_chain(request)
        .await
        .map(|res| res.into_inner());

    if response.is_ok() {
        create_timed_message(
            cx,
            "Chain deleted",
            format!("Chain {} was deleted", chain.name),
            InfoLevel::Success,
            5,
        );
        Ok(())
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
        Err(())
    }
}

pub async fn query_available_chains(cx: Scope<'_>) -> Result<Vec<Chain>, ()> {
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(());
    let response = services
        .grpc_client
        .get_settings_service()
        .get_available_chains(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(result) = response {
        Ok(result.chains)
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
        Err(())
    }
}

pub async fn add_chain(cx: Scope<'_>, chain: Chain) -> Result<Chain, ()> {
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(UpdateChainRequest {
        chain_id: chain.id,
        notify_new_proposals: true,
        notify_proposal_finished: true,
    });
    let response = services
        .grpc_client
        .get_settings_service()
        .add_chain(request)
        .await
        .map(|res| res.into_inner());

    if response.is_ok() {
        let new_chain = Chain {
            id: chain.id,
            name: chain.name,
            notify_new_proposals: true,
            notify_proposal_finished: true,
            logo_url: chain.logo_url,
            is_notify_new_proposals_supported: chain.is_notify_new_proposals_supported,
            is_notify_proposal_finished_supported: chain.is_notify_proposal_finished_supported,
        };
        Ok(new_chain)
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
        Err(())
    }
}
