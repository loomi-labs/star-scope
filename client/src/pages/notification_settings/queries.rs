use crate::components::messages::{
    create_error_msg_from_status, create_message, create_timed_message,
};
use crate::types::protobuf::grpc_settings::{
    RemoveWalletRequest, UpdateWalletRequest, ValidateWalletRequest, Wallet,
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
pub enum Update {
    Funding,
    Staking,
    GovVotingReminder,
}

pub async fn update_existing_wallet(cx: Scope<'_>, wallet_sig: &Signal<Wallet>, update: Update) {
    let wallet = create_ref(cx, wallet_sig.get_untracked());
    let notify_funding = if let Update::Funding = update {
        !wallet.notify_funding
    } else {
        wallet.notify_funding
    };
    let notify_staking = if let Update::Staking = update {
        !wallet.notify_staking
    } else {
        wallet.notify_staking
    };
    let notify_gov_voting_reminder = if let Update::GovVotingReminder = update {
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
            Update::Funding => {
                if notify_funding {
                    "You will be notified about funding events"
                } else {
                    "You will no longer be notified about funding events"
                }
            }
            Update::Staking => {
                if notify_staking {
                    "You will be notified about staking events"
                } else {
                    "You will no longer be notified about staking events"
                }
            }
            Update::GovVotingReminder => {
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
