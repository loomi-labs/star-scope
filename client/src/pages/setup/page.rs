use std::collections::{HashMap, HashSet};
use std::fmt::Display;
use std::hash::{Hash, Hasher};

use crate::components::button::{ColorScheme, OutlineButton, SolidButton};
use crate::components::search::{SearchEntity, Searchable};
use crate::config::keys;
use crate::pages::communication::page::{DiscordCard, TelegramCard};
use crate::utils::url::navigate_launch_app;
use sycamore::futures::spawn_local_scoped;
use sycamore::{prelude::*, view};
use tonic::Status;

use crate::components::loading::{LoadingSpinner, LoadingSpinnerSmall};
use crate::components::messages::{create_error_msg_from_status, create_message};
use crate::components::navigation::Header;
use crate::types::protobuf::grpc_user::step_response::Step;
use crate::types::protobuf::grpc_user::step_response::Step::{Five, Four, One, Three, Two};
use crate::types::protobuf::grpc_user::{
    finish_step_request, get_step_request, FinishStepRequest, GetStepRequest, GovChain,
    SearchWalletsRequest, StepFiveRequest, StepFiveResponse, StepFourRequest, StepFourResponse,
    StepOneRequest, StepResponse, StepThreeRequest, StepThreeResponse, StepTwoRequest,
    StepTwoResponse, User, ValidateWalletRequest, Validator, Wallet,
};
use crate::{AppRoutes, AppState, InfoLevel, Services};

#[derive(Debug, Clone)]
struct SetupState {
    pub step: RcSignal<Option<Step>>,
    pub num_steps: RcSignal<Option<usize>>,
}

impl SetupState {
    pub fn new() -> Self {
        Self {
            step: create_rc_signal(None),
            num_steps: create_rc_signal(None),
        }
    }
}

const TITLE_CLASS: &str = "text-3xl md:text-4xl font-bold m-4";
const SUBTITLE_CLASS: &str = "text-2xl font-semibold m-4";
const DESCRIPTION_CLASS: &str = "dark:text-purple-600 mb-8";
const DESCRIPTION_PROMINENT_CLASS: &str = "dark:text-white";
const BUTTON_ROW_CLASS: &str = "flex justify-center space-x-8 pb-6";

const SPACER_10_PERCENT_CLASS: &str = "w-full min-h-[10%]";

#[component]
fn StepOneComponent<G: Html>(cx: Scope) -> View<G> {
    let handle_click = move |is_validator| {
        spawn_local_scoped(cx, async move {
            let finish_step = FinishStepRequest {
                go_to_next_step: true,
                step: Some(finish_step_request::Step::One(StepOneRequest {
                    is_validator,
                })),
            };
            update_step(cx, finish_step).await;
        });
    };

    view! {cx,
        div() {
            div(class="hidden md:block") {
                lottie-player(src=keys::ROCKET_JSON, background="transparent", speed="1", style="width: 400px; height: 400px;", loop=true, autoplay=true) {}
            }
            div(class="md:hidden") {
                lottie-player(src=keys::ROCKET_JSON, background="transparent", speed="1", style="width: 300px; height: 300px;", loop=true, autoplay=true) {}
            }
        }
        div() {
            h1(class=TITLE_CLASS) {"Welcome to Star Scope!"}
            p(class=DESCRIPTION_CLASS) {"Setup your account to receive notifications about your Cosmos activities."}
        }
        div(class="mb-16") {
            h2(class=SUBTITLE_CLASS) {"Are you a validator?"}
            div(class="flex justify-center space-x-4") {
                SolidButton(on_click=move || handle_click(true), color=ColorScheme::Subtle) {"Yes"}
                SolidButton(on_click=move || handle_click(false), color=ColorScheme::Subtle) {"No"}
            }
        }
    }
}

#[component(inline_props)]
fn ProgressBar<G: Html>(cx: Scope, step: Step) -> View<G> {
    let setup_state: &SetupState = use_context::<SetupState>(cx);
    let num_steps_option = *setup_state.num_steps.get();
    if num_steps_option.is_none() {
        create_message(
            cx,
            "Invalid state",
            "Number of setup was not found",
            InfoLevel::Error,
        );
        return view! {cx,};
    }
    let num_steps = num_steps_option.unwrap() - 1;
    let current_step = match step {
        One(_) => 0,
        Two(_) => 0,
        Three(_) => num_steps - 3,
        Four(_) => num_steps - 2,
        Five(_) => num_steps - 1,
    };

    let current_name = match step {
        One(_) => "",
        Two(_) => "Validators",
        Three(_) => "Wallets",
        Four(_) => "Notifications",
        Five(_) => "Communication",
    };

    let progress_bar_views = View::new_fragment(
        (0..num_steps).map(|i|{
            let is_cicle_colored = i <= current_step;
            let has_line = i < num_steps - 1;
            let is_line_colored = current_step > i;
            let is_name_shown = i == current_step;
            view! { cx,
                (if is_cicle_colored {
                    view! {cx,
                        div(class="w-4 h-4 rounded-full bg-primary relative flex flex-col items-center") {
                          (if is_name_shown {
                            view! {cx,
                              p(class="absolute w-full flex-grow flex items-center justify-center -bottom-6 dark:text-purple-500") {
                                (current_name)
                              }
                            }
                          } else {
                            view! {cx,}
                          })
                        }
                      }
                } else {
                    view! {cx,
                        div(class="w-4 h-4 rounded-full dark:bg-purple-700")
                    }
                })
                (if has_line {
                    view! {cx,
                        (if is_line_colored {
                            view! {cx,
                                hr(class="w-1/8 flex-grow border border-primary")
                            }
                        } else {
                            view! {cx,
                                hr(class="w-1/8 flex-grow border dark:border-purple-700")
                            }
                        })
                    }
                } else {
                    view! {cx,}
                })
            }
        }).collect()
    );
    view! {cx,
        div(class=SPACER_10_PERCENT_CLASS)
        div(class="flex flex-col w-full justify-center items-center w-full space-y-2") {
            h2(class="text-xl dark:text-purple-600") { "Account setup" }
            div(class="flex justify-center items-center w-4/5 md:w-2/5") {
                (progress_bar_views)
            }
        }
        div(class=SPACER_10_PERCENT_CLASS)
    }
}

#[derive(Debug, Clone)]
struct SearchableValidator {
    pub validator: Validator,
}

impl Display for SearchableValidator {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        self.validator.moniker.fmt(f)
    }
}

impl Hash for SearchableValidator {
    fn hash<H: Hasher>(&self, state: &mut H) {
        self.validator.ids.hash(state);
    }
}

impl PartialEq for SearchableValidator {
    fn eq(&self, other: &Self) -> bool {
        self.validator.ids == other.validator.ids
    }
}

impl Eq for SearchableValidator {}

#[component(inline_props)]
fn StepTwoComponent<G: Html>(cx: Scope, step: StepTwoResponse) -> View<G> {
    let validators = step
        .available_validators
        .iter()
        .map(|val| {
            let is_selected = create_signal(
                cx,
                val.ids
                    .iter()
                    .all(|id| step.selected_validator_ids.contains(id)),
            );
            Searchable {
                entity: SearchableValidator {
                    validator: val.clone(),
                },
                is_selected,
            }
        })
        .collect::<Vec<Searchable<SearchableValidator>>>();

    let selected_validators = create_signal(cx, {
        validators
            .iter()
            .filter(|row| *row.is_selected.get_untracked())
            .map(|row| row.entity.clone())
            .collect::<HashSet<SearchableValidator>>()
    });

    let handle_click = move |go_to_next_step| {
        spawn_local_scoped(cx, async move {
            let finish_step = FinishStepRequest {
                go_to_next_step,
                step: Some(finish_step_request::Step::Two(StepTwoRequest {
                    validator_ids: selected_validators
                        .get()
                        .as_ref()
                        .clone()
                        .into_iter()
                        .flat_map(|val| val.validator.ids)
                        .collect(),
                })),
            };
            update_step(cx, finish_step).await;
        });
    };

    view! {cx,
        ProgressBar(step=Two(step))
        div(class="flex flex-col md:h-2/5 md:w-3/5") {
            h2(class=TITLE_CLASS) {"Choose your validator(s)"}
            p(class=DESCRIPTION_CLASS) {"You will receive reminders to vote on governance proposals from the validators you've selected."}
            SearchEntity(searchables=validators.clone(), selected_entities=selected_validators.clone(), placeholder="Search validators")
        }
        div(class="min-h-[10%] md:min-h-0")
        div(class=BUTTON_ROW_CLASS) {
            OutlineButton(on_click=move || handle_click(false)) {"Back"}
            SolidButton(on_click=move || handle_click(true)) {"Next"}
        }
        div(class=SPACER_10_PERCENT_CLASS)
    }
}

#[derive(Debug, Clone, PartialEq)]
enum WalletValidation {
    Valid(Wallet),
    Invalid(String),
}

async fn query_validate_wallet(cx: Scope<'_>, address: String) -> WalletValidation {
    let services = use_context::<Services>(cx);
    let request = services
        .grpc_client
        .create_request(ValidateWalletRequest { address });
    let response = services
        .grpc_client
        .get_user_setup_service()
        .validate_wallet(request)
        .await
        .map(|res| res.into_inner());
    if let Ok(response) = response {
        if response.is_valid {
            if response.is_supported {
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

#[derive(Debug, Clone, PartialEq)]
struct NewWallet<'a> {
    wallet: Wallet,
    parent_wallet_address: Option<String>,
    is_added: &'a Signal<bool>,
}

#[component(inline_props)]
fn AddWallet<'a, G: Html>(cx: Scope<'a>, wallets: &'a Signal<Vec<NewWallet<'a>>>) -> View<G> {
    let new_wallet_address = create_signal(cx, String::new());

    let validation = create_signal(cx, None::<WalletValidation>);

    create_effect(cx, move || {
        let address = new_wallet_address.get().as_ref().clone();
        if wallets
            .get()
            .iter()
            .map(|w| w.wallet.address.clone())
            .collect::<Vec<String>>()
            .contains(&new_wallet_address.get().as_ref().clone())
        {
            validation.set(Some(WalletValidation::Invalid(
                "Wallet already added".to_string(),
            )));
        } else if address.is_empty() {
            validation.set(None);
        } else if address.len() < 30 {
            validation.set(Some(WalletValidation::Invalid(
                "Wallet address is too short".to_string(),
            )));
        } else {
            spawn_local_scoped(cx, async move {
                let result =
                    query_validate_wallet(cx, new_wallet_address.get().as_ref().clone()).await;
                validation.set(Some(result.clone()));
                if let WalletValidation::Valid(wallet) = result {
                    wallets.modify().push(NewWallet {
                        wallet,
                        parent_wallet_address: None,
                        is_added: create_signal(cx, true),
                    });
                    new_wallet_address.set(String::new());
                }
            });
        }
    });

    view! {cx,
        div(class="flex flex-col") {
            div(class="flex") {
                input(
                    class="w-full placeholder:italic border border-gray-300 rounded-lg px-4 py-2 text-black focus:outline-none focus:ring-2 focus:ring-primary",
                    placeholder="Wallet address",
                    type="text",
                    bind:value=new_wallet_address,
                )
            }
            (if let Some(WalletValidation::Invalid(msg)) = validation.get().as_ref().clone() {
                view! {cx,
                    span(class="text-red-500 text-left") {(msg)}
                }
            } else {
                view! {cx, }
            })
        }
    }
}

async fn query_search_wallets<'a>(
    cx: Scope<'a>,
    address: String,
    wallets: &Signal<Vec<NewWallet<'a>>>,
    searched: Vec<String>,
) {
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(SearchWalletsRequest {
        address: address.clone(),
        added_addresses: wallets
            .get()
            .iter()
            .map(|w| w.wallet.address.clone())
            .collect(),
        searched_bech32_addresses: searched,
    });
    let result = services
        .grpc_client
        .get_user_setup_service()
        .search_wallets(request)
        .await
        .map(|res| res.into_inner());
    match result {
        Ok(mut stream) => loop {
            match stream.message().await {
                Ok(Some(response)) => {
                    let setup_state = try_use_context::<SetupState>(cx);
                    if setup_state.is_none() {
                        break;
                    }
                    if !matches!(*setup_state.unwrap().step.get_untracked(), Some(Three(_))) {
                        break;
                    }
                    if response.wallet.is_none() {
                        create_message(cx, "Error", "No wallets found", InfoLevel::Error);
                    } else {
                        let is_already_added = wallets
                            .get()
                            .iter()
                            .map(|w| w.wallet.address.clone())
                            .collect::<Vec<String>>()
                            .contains(&response.wallet.as_ref().unwrap().address.clone());
                        if !is_already_added {
                            wallets.modify().push(NewWallet {
                                wallet: response.wallet.unwrap(),
                                parent_wallet_address: Some(address.clone()),
                                is_added: create_signal(cx, false),
                            });
                        }
                    }
                }
                Ok(None) => break,
                Err(err) => create_error_msg_from_status(cx, err),
            }
        },
        Err(err) => {
            create_error_msg_from_status(cx, err);
        }
    }
}

#[derive(Debug, Clone, PartialEq)]
struct SearchQuery<'a> {
    bech32_address: String,
    is_searching: &'a Signal<bool>,
}

#[component(inline_props)]
fn WalletList<'a, G: Html>(cx: Scope<'a>, wallets: &'a Signal<Vec<NewWallet<'a>>>) -> View<G> {
    let search_wallets = create_signal(cx, Vec::<SearchQuery>::new());

    create_effect(cx, move || {
        let searched_addresses = search_wallets
            .get_untracked()
            .iter()
            .cloned()
            .map(|w| w.bech32_address)
            .collect::<HashSet<String>>();

        let unique_addresses: HashMap<String, NewWallet> = wallets
            .get()
            .iter()
            .filter(|w| !searched_addresses.contains(&w.wallet.bech32_address))
            .map(|w| (w.wallet.bech32_address.clone(), w.clone()))
            .collect();

        unique_addresses
            .into_iter()
            .for_each(|(bech32_address, w)| {
                let search_wallet: SearchQuery<'_> = SearchQuery {
                    bech32_address,
                    is_searching: create_signal(cx, true),
                };
                let cloned_search_wallet = search_wallet.clone();
                search_wallets.modify().push(search_wallet);
                let searched = searched_addresses.iter().cloned().collect::<Vec<String>>();
                spawn_local_scoped(cx, async move {
                    query_search_wallets(cx, w.wallet.address.clone(), wallets, searched).await;
                    cloned_search_wallet.is_searching.set(false);
                });
            });
    });

    let is_searching = create_selector(cx, {
        let search_wallets = search_wallets.clone();
        move || search_wallets.get().iter().any(|w| *w.is_searching.get())
    });

    let cnt_wallets = create_signal(cx, wallets.get().len());
    let cnt_non_added_wallets = create_signal(
        cx,
        wallets
            .get()
            .iter()
            .filter(|w| !w.is_added.get().as_ref())
            .count(),
    );

    create_effect(cx, move || {
        if *is_searching.get() {
            return;
        }
        let old_cnt = *cnt_wallets.get_untracked();
        let new_cnt = wallets.get_untracked().len();
        let old_non_added_cnt = *cnt_non_added_wallets.get_untracked();
        let new_non_added_cnt = wallets
            .get_untracked()
            .iter()
            .filter(|w| !w.is_added.get().as_ref())
            .count();
        cnt_wallets.set(new_cnt);
        cnt_non_added_wallets.set(new_non_added_cnt);
        if old_cnt < new_cnt && old_non_added_cnt != new_non_added_cnt {
            create_message(
                cx,
                "Wallets found",
                format!("{} wallets found", new_cnt - old_cnt),
                InfoLevel::Success,
            );
        }
    });

    view! {cx,
        div(class="flex flex-col space-y-4") {
            Indexed(
                iterable=wallets,
                view=move |cx, wallet| {
                    let wallet_ref = create_ref(cx, wallet);
                    let prefix = wallet_ref.wallet.address[..8].to_owned();
                    let suffix = wallet_ref.wallet.address[wallet_ref.wallet.address.len() - 4..].to_owned();
                    let shortened_address = format!("{}...{}", prefix, suffix);
                    view! {cx,
                        div(class="flex justify-between items-center space-x-8") {
                            div(class=format!("flex flex-grow items-center justify-center space-x-2 px-4 py-2 rounded-full dark:bg-purple-700 {}", if *wallet_ref.is_added.get() { "" } else { "opacity-50" })) {
                                img(src=wallet_ref.wallet.logo_url, alt="Chain logo", class="h-6 w-6")
                                span(class="text-sm") {(shortened_address)}
                            }
                            (if *wallet_ref.is_added.get() {
                                view!{cx,
                                    button(class="flex justify-between items-center p-2 rounded-lg border-2 border-purple-700 hover:bg-primary",
                                            on:click=move |_| wallet_ref.is_added.set(false)) {
                                        span(class="w-6 h-6 icon-[tabler--trash] cursor-pointer")
                                    }
                                }
                            } else {
                                view!{cx,
                                    button(class="flex justify-between items-center p-2 opacity-100 rounded-lg border-2 border-green-500 bg-green-500 hover:bg-green-600 hover:border-green-600",
                                            on:click=move |_| wallet_ref.is_added.set(true)) {
                                        span(class="w-6 h-6 icon-[ic--round-add] cursor-pointer")
                                    }
                                }
                            })
                        }
                    }
                }
            )
            (if *is_searching.get() {
                view! {cx,
                    LoadingSpinnerSmall {}
                }
            } else {
                view! {cx, }
            })
            AddWallet(wallets=wallets.clone())
        }
    }
}

#[component(inline_props)]
fn StepThreeComponent<G: Html>(cx: Scope, step: StepThreeResponse) -> View<G> {
    let wallets: &Signal<Vec<NewWallet>> = create_signal(
        cx,
        step.wallets
            .iter()
            .cloned()
            .map(|wallet| NewWallet {
                wallet,
                parent_wallet_address: None,
                is_added: create_signal(cx, true),
            })
            .collect(),
    );

    let handle_click = move |go_to_next_step| {
        spawn_local_scoped(cx, async move {
            let finish_step = FinishStepRequest {
                go_to_next_step,
                step: Some(finish_step_request::Step::Three(StepThreeRequest {
                    wallet_addresses: wallets
                        .get()
                        .iter()
                        .filter(|w| *w.is_added.get())
                        .map(|w| w.wallet.address.clone())
                        .collect(),
                })),
            };
            update_step(cx, finish_step).await;
        });
    };

    view! {cx,
        ProgressBar(step=Three(step))
        h2(class=TITLE_CLASS) {"Add your wallets"}
        p(class=DESCRIPTION_CLASS) {"You will receive notifications about events related to your wallets."}
        WalletList(wallets=wallets.clone())
        div(class=SPACER_10_PERCENT_CLASS)
        div(class=BUTTON_ROW_CLASS) {
            OutlineButton(on_click=move || handle_click(false)) {"Back"}
            SolidButton(on_click=move || handle_click(true)) {"Next"}
        }
    }
}

#[derive(Debug, Clone)]
struct SearchableChain {
    pub chain: GovChain,
}

impl Display for SearchableChain {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        self.chain.name.fmt(f)
    }
}

impl Hash for SearchableChain {
    fn hash<H: Hasher>(&self, state: &mut H) {
        self.chain.id.hash(state);
    }
}

impl PartialEq for SearchableChain {
    fn eq(&self, other: &Self) -> bool {
        self.chain.id == other.chain.id
    }
}

impl Eq for SearchableChain {}

#[component(inline_props)]
fn StepFourComponent<G: Html>(cx: Scope, step: StepFourResponse) -> View<G> {
    let notify_funding = create_signal(cx, step.notify_funding);
    let notify_staking = create_signal(cx, step.notify_staking);
    let notify_gov_new_proposal = create_signal(cx, step.notify_gov_new_proposal);
    let notify_gov_voting_end = create_signal(cx, step.notify_gov_voting_end);
    let notify_gov_voting_reminder = create_signal(cx, step.notify_gov_voting_reminder);

    let chain_rows = step
        .available_chains
        .iter()
        .map(|chain| {
            let is_selected = create_signal(cx, step.notify_gov_chain_ids.contains(&chain.id));
            Searchable {
                entity: SearchableChain {
                    chain: chain.clone(),
                },
                is_selected,
            }
        })
        .collect::<Vec<Searchable<SearchableChain>>>();

    let selected_chains = create_signal(cx, {
        chain_rows
            .iter()
            .filter(|row| *row.is_selected.get_untracked())
            .map(|row| row.entity.clone())
            .collect::<HashSet<SearchableChain>>()
    });

    let handle_click = move |go_to_next_step| {
        spawn_local_scoped(cx, async move {
            let finish_step = FinishStepRequest {
                go_to_next_step,
                step: Some(finish_step_request::Step::Four(StepFourRequest {
                    notify_funding: *notify_funding.get(),
                    notify_staking: *notify_staking.get(),
                    notify_gov_new_proposal: *notify_gov_new_proposal.get(),
                    notify_gov_voting_end: *notify_gov_voting_end.get(),
                    notify_gov_voting_reminder: *notify_gov_voting_reminder.get(),
                    notify_gov_chain_ids: selected_chains
                        .get()
                        .as_ref()
                        .clone()
                        .into_iter()
                        .map(|chain| chain.chain.id)
                        .collect(),
                })),
            };
            update_step(cx, finish_step).await;
        });
    };

    let section_selected_class = "w-6 h-6 bg-primary icon-[icon-park-solid--check-one]";
    let section_unselected_class = "w-6 h-6 rounded-full border-2 border-primary";
    let centered_row_class = "flex justify-center items-center space-x-4";
    let starting_row_class = "flex items-center space-x-4 p-2";
    let starting_row_selectable_class = create_ref(
        cx,
        format!("{} rounded-lg hover:dark:bg-purple-700", starting_row_class),
    );
    let check_mark_class = "w-6 h-6 bg-primary icon-[ph--check-bold]";

    view! {cx,
        ProgressBar(step=Four(step))
        h2(class=TITLE_CLASS) {"Choose your notifications"}
        div(class="flex flex-wrap w-full lg:w-4/5 rounded-xl mt-4 text-sm md:text-base dark:bg-purple-800") {
            div(class="flex flex-col items-center w-full md:w-1/3 px-6 py-10 md:py-6 rounded-xl hover:dark:bg-purple-700", on:click=move |_| notify_funding.set(!notify_funding.get().as_ref())) {
                div(class=centered_row_class) {
                    span(class=(if *notify_funding.get() {section_selected_class} else {section_unselected_class})) {}
                    h3(class=SUBTITLE_CLASS) {"Funding"}
                }
                p(class=DESCRIPTION_PROMINENT_CLASS) {"When you receive tokens"}
            }
            div(class="flex w-full md:w-1/3 rounded-xl hover:dark:bg-purple-700", on:click=move |_| notify_staking.set(!notify_staking.get().as_ref())) {
                div(class="flex flex-col md:flex-row w-full h-full") {
                    div(class="flex flex-col md:flex-row items-center w-full md:w-0") {
                        div(class="border border-purple-500 w-4/5 md:w-0 md:h-4/5") {}
                    }
                    div(class="flex flex-col flex-grow items-center p-6") {
                        div(class=centered_row_class) {
                            span(class=(if *notify_staking.get() {section_selected_class} else {section_unselected_class})) {}
                            h3(class=SUBTITLE_CLASS) {"Staking"}
                        }
                        div(class="flex flex-col") {
                            div(class=starting_row_class) {
                                span(class=check_mark_class)
                                span() {"When unbonding period is over"}
                            }
                            div(class=starting_row_class) {
                                span(class=check_mark_class)
                                span() {"When a validator gets slashed"}
                            }
                            div(class=starting_row_class) {
                                span(class=check_mark_class)
                                span() {"When a validator gets inactive"}
                            }
                        }
                    }
                    div(class="flex flex-col md:flex-row items-center w-full md:w-0") {
                        div(class="border border-purple-500 w-4/5 md:w-0 md:h-4/5") {}
                    }
                }
            }
            div(class="flex flex-col items-center w-full md:w-1/3 p-2 rounded-xl p-6") {
                h3(class=SUBTITLE_CLASS) {"Governance"}
                div(class="flex flex-col mb-4") {
                    div(class=starting_row_selectable_class, on:click=move |_| notify_gov_new_proposal.set(!notify_gov_new_proposal.get().as_ref())) {
                        span(class=(if *notify_gov_new_proposal.get() {section_selected_class} else {section_unselected_class})) {}
                        span() {"New proposal open for voting"}
                    }
                    div(class=starting_row_selectable_class, on:click=move |_| notify_gov_voting_end.set(!notify_gov_voting_end.get().as_ref())) {
                        span(class=(if *notify_gov_voting_end.get() {section_selected_class} else {section_unselected_class})) {}
                        span() {"Proposal passed/failed"}
                    }
                    div(class=starting_row_selectable_class, on:click=move |_| notify_gov_voting_reminder.set(!notify_gov_voting_reminder.get().as_ref())) {
                        span(class=(if *notify_gov_voting_reminder.get() {section_selected_class} else {section_unselected_class})) {}
                        span() {"Voting reminders"}
                    }
                }
                SearchEntity(searchables=chain_rows, selected_entities=selected_chains, placeholder="Select chains", show_results_for_empty_search=true)
            }
        }
        div(class=SPACER_10_PERCENT_CLASS)
        div(class=BUTTON_ROW_CLASS) {
            OutlineButton(on_click=move || handle_click(false)) {"Back"}
            SolidButton(on_click=move || handle_click(true)) {"Next"}
        }
    }
}

#[component(inline_props)]
fn StepFiveComponent<G: Html>(cx: Scope, step: StepFiveResponse) -> View<G> {
    let handle_click = move |go_to_next_step| {
        spawn_local_scoped(cx, async move {
            let finish_step = FinishStepRequest {
                go_to_next_step,
                step: Some(finish_step_request::Step::Five(StepFiveRequest {
                    ..Default::default()
                })),
            };
            update_step(cx, finish_step).await;
        });
    };
    let web_app_url = create_ref(
        cx,
        keys::WEB_APP_URL.to_string() + AppRoutes::Setup.to_string().as_str(),
    );

    view! {cx,
        ProgressBar(step=Five(step))
        h2(class=TITLE_CLASS) {"Choose your notification channels"}
        div(class=DESCRIPTION_CLASS) {
            "You will always receive notifications on the web app."
            br()
            "Choose additional channels to receive notifications on."
        }
        div(class="flex flex-col space-x-0 space-y-8 mt-4 md:flex-row md:space-x-8 md:space-y-0") {
            div(class="flex items-center p-8 rounded-lg md:w-1/2 dark:bg-purple-700") {
                DiscordCard(web_app_url=web_app_url.clone(), center_button=true)
            }
            div(class="p-8 rounded-lg md:w-1/2 dark:bg-purple-700") {
                TelegramCard(web_app_url=web_app_url.clone(), center_button=true)
            }
        }
        div(class=SPACER_10_PERCENT_CLASS)
        div(class=BUTTON_ROW_CLASS) {
            OutlineButton(on_click=move || handle_click(false)) {"Back"}
            SolidButton(on_click=move || handle_click(true)) {"Finish"}
        }
    }
}

async fn update_step(cx: Scope<'_>, finish_step: FinishStepRequest) {
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(finish_step);
    let response = services
        .grpc_client
        .get_user_setup_service()
        .finish_step(request)
        .await
        .map(|res| res.into_inner());
    handle_setup_step_response(cx, response)
}

fn handle_setup_step_response(cx: Scope, response: Result<StepResponse, Status>) {
    if let Ok(response) = response {
        let setup_state = use_context::<SetupState>(cx);
        match response.step {
            Some(step) => {
                match response.num_steps.try_into() {
                    Ok(num_steps) => setup_state.num_steps.set(Some(num_steps)),
                    Err(err) => create_message(cx, "Error", err.to_string(), InfoLevel::Error),
                }
                if response.is_complete {
                    let app_state = use_context::<AppState>(cx);
                    if let Some(user) = app_state.user.get_untracked().as_ref().clone() {
                        let new_user = User {
                            is_setup_complete: true,
                            ..user
                        };
                        app_state.user.set(Some(new_user));
                        navigate_launch_app(cx)
                    } else {
                        create_message(
                            cx,
                            "User not found",
                            "User status unknown",
                            InfoLevel::Error,
                        );
                    }
                }
                setup_state.step.set(Some(step));
            }
            None => create_message(
                cx,
                "Setup step not found",
                "Setup step was not found",
                crate::InfoLevel::Error,
            ),
        }
    } else {
        create_error_msg_from_status(cx, response.err().unwrap());
    }
}

async fn query_setup_step(cx: Scope<'_>, requestedStep: Option<get_step_request::Step>) {
    let services = use_context::<Services>(cx);
    let request = services.grpc_client.create_request(GetStepRequest {
        step: requestedStep.unwrap_or(get_step_request::Step::Current) as i32,
    });
    let response = services
        .grpc_client
        .get_user_setup_service()
        .get_step(request)
        .await
        .map(|res| res.into_inner());
    handle_setup_step_response(cx, response)
}

#[component]
pub fn Setup<G: Html>(cx: Scope) -> View<G> {
    let setup_state = SetupState::new();
    provide_context(cx, setup_state.clone());

    spawn_local_scoped(cx, async move {
        query_setup_step(cx, None).await;
    });

    view! {cx,
        (if let Some(step) = setup_state.step.get().as_ref() {
            let child = match step {
                One(_) => view! {cx, StepOneComponent()},
                Two(s) => view! {cx, StepTwoComponent(step=s.clone())},
                Three(s) => view! {cx, StepThreeComponent(step=s.clone())},
                Four(s) => view! {cx, StepFourComponent(step=s.clone())},
                Five(s) => view! {cx, StepFiveComponent(step=s.clone())},
            };
            view! {
                cx,
                div(class="h-[100dvh] flex justify-center items-center flex-auto flex-shrink-0") {
                    div(class="w-full h-full flex flex-col lg:max-w-screen-lg xl:max-w-screen-xl overflow-y-auto") {
                        Header{}
                        div(class="w-full h-full flex flex-col flex-auto items-center text-center p-4") {
                            (child)
                        }
                    }
                }
            }
        } else {
            view! {cx,
                LoadingSpinner {}
            }
        })
    }
}
