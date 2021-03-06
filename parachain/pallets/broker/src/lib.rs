#![cfg_attr(not(feature = "std"), no_std)]
#![allow(unused_variables)]

use frame_support::{decl_module, decl_storage, decl_event, decl_error, dispatch::DispatchResult};
use frame_system::{self as system, ensure_signed, ensure_root};

use sp_std::prelude::*;

use common::{AppID, Message, Verifier, Broker, Application, registry::{REGISTRY, AppName}};

pub trait Trait: system::Trait {

	type Event: From<Event> + Into<<Self as system::Trait>::Event>;

	type DummyVerifier: Verifier;
	type PolkaETH: Application;

}

decl_storage! {
	trait Store for Module<T: Trait> as BrokerModule {

	}
}

decl_event!(
	pub enum Event {
	}
);

decl_error! {
	pub enum Error for Module<T: Trait> {
		HandlerNotFound
	}
}

decl_module! {
	pub struct Module<T: Trait> for enum Call where origin: T::Origin {

		type Error = Error<T>;

		fn deposit_event() = default;

		// Accept a message that has been previously verified
		//
		// Called by a verifier
		#[weight = 0]
		pub fn accept(origin, app_id: AppID, message: Message) -> DispatchResult {
			// we'll check the origin here to ensure it originates from a verifier
			//ensure_root(origin)?;
			Self::dispatch(app_id, message)?;
			Ok(())
		}

	}
}

impl<T: Trait> Module<T> {

	// Dispatch verified message to a target application
	fn dispatch(app_id: AppID, message: Message) -> DispatchResult {

		for entry in REGISTRY.iter() {
			if app_id != entry.id {
				continue
			}
			return match entry.symbol {
				AppName::PolkaETH => {
					T::PolkaETH::handle(app_id, message)?;
					Ok(())
				}
			}
		}

		Err(<Error<T>>::HandlerNotFound.into())
	}
}


impl<T: Trait> Broker for Module<T> {

	// Submit message to broker for processing
	//
	// The message will be routed to a verifier
	fn submit(app_id: AppID, message: Message) -> DispatchResult {
		T::DummyVerifier::verify(app_id, message)?;
		Ok(())
	}

}
