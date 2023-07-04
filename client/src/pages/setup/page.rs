use std::collections::HashSet;
use std::fmt::Display;
use std::hash::{Hash, Hasher};

use crate::components::button::{ColorScheme, OutlineButton, SolidButton};
use crate::components::search::{SearchEntity, Searchable};
use crate::utils::url::navigate_launch_app;
use sycamore::futures::spawn_local_scoped;
use sycamore::{prelude::*, view};
use tonic::Status;

use crate::components::loading::LoadingSpinner;
use crate::components::messages::{create_error_msg_from_status, create_message};
use crate::components::navigation::Header;
use crate::types::protobuf::grpc::step_response::Step;
use crate::types::protobuf::grpc::step_response::Step::{
    Five, Four, One, Three, Two,
};
use crate::types::protobuf::grpc::{
    finish_step_request, get_step_request, FinishStepRequest, GetStepRequest, GovChain,
    StepFiveRequest, StepFiveResponse, StepFourRequest, StepFourResponse, StepOneRequest,
    StepResponse, StepThreeRequest, StepThreeResponse, StepTwoRequest, StepTwoResponse, User,
    ValidateWalletRequest, Validator, Wallet,
};
use crate::{AppState, InfoLevel, Services};

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

const TITLE_CLASS: &str = "text-4xl font-bold m-4";
const SUBTITLE_CLASS: &str = "text-2xl font-semibold m-4";
const DESCRIPTION_CLASS: &str = "dark:text-purple-600 mb-8";
const DESCRIPTION_PROMINENT_CLASS: &str = "dark:text-white";
const BUTTON_ROW_CLASS: &str = "flex justify-center space-x-8 mt-8";

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
        h1(class=TITLE_CLASS) {"Welcome to Star Scope!"}
        p(class=DESCRIPTION_CLASS) {"We deliver quick and effortless notifications about your Cosmos ecosystem activities."}
        h2(class=SUBTITLE_CLASS) {"Are you a validator?"}
        div(class="flex justify-center space-x-4") {
            SolidButton(on_click=move || handle_click(true), color=ColorScheme::Subtle) {"Yes"}
            SolidButton(on_click=move || handle_click(false), color=ColorScheme::Subtle) {"No"}
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
        div(class="flex flex-col justify-center items-center w-full space-y-2") {
            h2(class="text-xl dark:text-purple-600") { "Account setup" }
            div(class="flex justify-center items-center w-4/5 md:w-2/5") {
                (progress_bar_views)
            }
            div(class="w-full h-10")
        }
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
        h2(class=TITLE_CLASS) {"Choose your validator(s)"}
        p(class=DESCRIPTION_CLASS) {"You will receive reminders to vote on governance proposals from the validators you've selected."}
        SearchEntity(searchables=validators.clone(), selected_entities=selected_validators.clone(), placeholder="Search validators")
        div(class=BUTTON_ROW_CLASS) {
            OutlineButton(on_click=move || handle_click(false)) {"Back"}
            SolidButton(on_click=move || handle_click(true)) {"Next"}
        }
    }
}

#[component(inline_props)]
fn WalletList<'a, G: Html>(cx: Scope<'a>, wallets: &'a Signal<Vec<Wallet>>) -> View<G> {
    let handle_delete_wallet = move |wallet: &Wallet| {
        wallets.modify().retain(|w| w.address != wallet.address);
    };

    view! {cx,
        div(class="flex flex-col space-y-4") {
            Indexed(
                iterable=wallets,
                view=move |cx, wallet| {
                    let wallet_ref = create_ref(cx, wallet);
                    let prefix = wallet_ref.address[..8].to_owned();
                    let suffix = wallet_ref.address[wallet_ref.address.len() - 4..].to_owned();
                    let shortened_address = format!("{}...{}", prefix, suffix);
                    view! {cx,
                        div(class="flex justify-between items-center space-x-8") {
                            div(class="flex flex-grow items-center justify-center space-x-2 px-4 py-2 rounded-full dark:bg-purple-700") {
                                img(src=wallet_ref.logo_url, alt="Chain logo", class="h-6 w-6")
                                span(class="text-sm") {(shortened_address)}
                            }
                            button(class="flex justify-between items-center p-2 rounded-lg border-2 border-purple-700 hover:bg-primary",
                                    on:click=move |_| handle_delete_wallet(wallet_ref)) {
                                span(class="w-6 h-6 icon-[tabler--trash] cursor-pointer")
                            }
                        }
                    }
                }
            )
            AddWallet(wallets=wallets.clone())
        }
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

#[component(inline_props)]
fn AddWallet<'a, G: Html>(cx: Scope<'a>, wallets: &'a Signal<Vec<Wallet>>) -> View<G> {
    let new_wallet_address = create_signal(cx, String::new());

    let validation = create_signal(cx, None::<WalletValidation>);

    create_effect(cx, move || {
        let address = new_wallet_address.get().as_ref().clone();
        if wallets
            .get()
            .iter()
            .map(|w| w.address.clone())
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
                    wallets.modify().push(wallet);
                    new_wallet_address.set(String::new());
                }
            });
        }
    });

    view! {cx,
        div(class="flex flex-col") {
            div(class="flex") {
                input(
                    class="w-full border border-gray-300 rounded-lg px-4 py-2 text-black focus:outline-none focus:ring-2 focus:ring-primary",
                    placeholder="Wallet address",
                    type="text",
                    bind:value=new_wallet_address,
                )
                // SolidButton(color=ColorScheme::Subtle, on_click=handle_add_wallet) {"Add"}
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

#[component(inline_props)]
fn StepThreeComponent<G: Html>(cx: Scope, step: StepThreeResponse) -> View<G> {
    let wallets = create_signal(cx, step.wallets.clone());

    let handle_click = move |go_to_next_step| {
        spawn_local_scoped(cx, async move {
            let finish_step = FinishStepRequest {
                go_to_next_step,
                step: Some(finish_step_request::Step::Three(StepThreeRequest {
                    wallet_addresses: wallets.get().iter().map(|w| w.address.clone()).collect(),
                })),
            };
            update_step(cx, finish_step).await;
        });
    };

    view! {cx,
        ProgressBar(step=Three(step))
        h2(class=TITLE_CLASS) {"Add your wallet(s)"}
        p(class=DESCRIPTION_CLASS) {"You will receive notifications about important updates and events directly related to your wallet."}
        WalletList(wallets=wallets.clone())
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
        div(class="flex flex-wrap rounded-xl mt-4 text-sm md:text-base dark:bg-purple-800") {
            div(class="flex flex-col items-center w-full md:w-1/3 px-6 py-10 md:py-6 rounded-xl hover:dark:bg-purple-700", on:click=move |_| notify_funding.set(!notify_funding.get().as_ref())) {
                div(class=centered_row_class) {
                    span(class=(if *notify_funding.get() {section_selected_class} else {section_unselected_class})) {}
                    h3(class=SUBTITLE_CLASS) {"Funding"}
                }
                p(class=DESCRIPTION_PROMINENT_CLASS) {"Whenever someone sends you tokens, we'll make sure you know"}
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
                                span() {"When a validator falls out of the active set"}
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
                        span() {"New proposal in voting period"}
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
                SearchEntity(searchables=chain_rows, selected_entities=selected_chains, placeholder="Select chains")
            }
        }
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

    view! {cx,
        ProgressBar(step=Five(step))
        p(class=DESCRIPTION_CLASS) {"Choose your notification channel"}
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
                div(class="min-h-[100svh] flex justify-center items-center flex-auto flex-shrink-0") {
                    div(class="min-h-[100svh] flex flex-col lg:max-w-screen-lg xl:max-w-screen-xl h-full w-full") {
                        Header{}
                        div(class="w-full h-full flex flex-col flex-auto justify-center text-center p-4") {
                            div(class="flex flex-col justify-center items-center") {
                                (child)
                            }
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
