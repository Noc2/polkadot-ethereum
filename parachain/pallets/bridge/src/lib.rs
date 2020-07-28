#![cfg_attr(not(feature = "std"), no_std)]

use frame_support::{decl_error, decl_event, decl_module, decl_storage, dispatch};
use frame_system::{self as system, ensure_signed};

use sp_std::prelude::*;

use sp_runtime::traits::Hash;

use common::{debug, AppID, Bridge, Broker, Message};

pub trait Trait: system::Trait {
	type Event: From<Event<Self>> + Into<<Self as system::Trait>::Event>;

	type Broker: Broker;
}

decl_storage! {

	trait Store for Module<T: Trait> as BridgeModule {
	}
}

decl_event!(
	pub enum Event<T>
	where
		AccountId = <T as system::Trait>::AccountId,
		Hash = <T as frame_system::Trait>::Hash,
	{
		MessageReceived(AccountId, AppID, Hash),
		AppEvent(AppID, Vec<u8>, Vec<u8>),
	}
);

decl_error! {
	pub enum Error for Module<T: Trait> {
	}
}

decl_module! {

	pub struct Module<T: Trait> for enum Call where origin: T::Origin {

		type Error = Error<T>;

		fn deposit_event() = default;

		#[weight = 0]
		pub fn send(origin, app_id: AppID, message: Message) -> dispatch::DispatchResult {
			debug!("foo {:?}", app_id);
			let who = ensure_signed(origin)?;
			Self::deposit_event(RawEvent::MessageReceived(who, app_id, T::Hashing::hash(message.as_ref())));

			T::Broker::submit(app_id, message.clone())?;
			Ok(())
		}

	}
}

impl<T: Trait> Bridge for Module<T> {
	fn deposit_event(app_id: AppID, name: Vec<u8>, data: Vec<u8>) {
		Self::deposit_event(RawEvent::AppEvent(app_id, name, data));
	}
}
