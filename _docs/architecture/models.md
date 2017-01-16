# Data models
To appropriately describe features the data we store has to be defined first. This 
data will be stored on the server and then retrieved by the client via an api. 

This data is organized in a flat relational structure so it can be easily 
serialized, cached, and sent to clients.

## Metadata
Every single model contains the following fields which describe some metadata:

- `id`: `int`
    - A unique id so the object can be referenced
- `created_at`: `Time`
    - Time the model was created.
- `updated_at`: `Time`
    - Time the model was last updated.
- `deleted_at`: `Time`
    - Time the model was delete, if not delete than null.
    
# User
A Squad Up user. Only one user per email.

## Model
### User
Information about user.

- `first_name`: `string`
- `last_name`: `string`
- `email`: `string`
- `profile_picture_url`: `string`

# Location
Locations can be owned by users or squads. Identified by (lat, long). App should offer some helper to 
resolve addresses into lat long format.

## Model
### Location
- `name`: `string`
- `owner_type`: `string`
    - Either `user` or `squad`
- `owner_id`: `int`
    - Id of user if `owner_type=user` and id of squad if `owner_type=squad`
- `lat`: `long`
    - Latitude of location
- `long`: `long`
    - Longitude of location

# Squad
A group of users.

## Model
### Squad
Represents a squad

- `name`: `string`
- `logo_url`: `string`

### Squad Membership
Tracks membership to squads.

- `squad_id`: `int`
    - Id of squad membership refers to
- `user_id`: `int`
    - Id of user membership refers to
    
# Car
Describes a car, users own cars.

## Model
- `user_id`: `int`
    - Id of user who owns car
- `seats`: `int`
- `make`: `string`
    - Make of the car
- `model`: `string`
    - Model of the car
- `color`: `string`
- `loc_lat`: `long`
    - If current location is known: Location latitude, if unknown then `-1`.
- `loc_long`: `long`
    - If current location is known: Location longitude, if unknown then `-1`.

# Event
Squad event. An event holds basic info about an event and then some optional extra contextual info 
in this format:

## Model
### Event
General event information

- `squad_id`: `int`
- `name`: `string`
- `description`: `string`
- `organizer_id`: `int`
    - Id of user who created event.
- `logo_url`: `string`
- `target_location_lat`: `long`
    - See param below, `target_loc_rad`
- `target_loc_long`: `long`
    - See param below, `target_loc_rad`
- `target_loc_rad`: `long`
    - The combination of the fields `lat`, `long`, and `rad` (Radius) (Prefixed with `target_loc`) 
      represents an area (circle) on a map that the event is taking place.
    - This acts as an anchor for all the algorithms when trying to determine the best locations 
      as a source for the filters (see Event technical docs for more detail)
- `final_event_proposal_id`: `int`
    - The id of the final decided on event proposal or `-1` if undecided.
- `type`: `string`
    - Describes what general type of event it is. The general types are:
        - `Get Together`
            - Typical behavior
                - Users all meet at one location.
                - Stay at location from `start_time` to `end_time`.
            - Meta option: None
            - Sources:
                - Proposed Location
                    - Given by squad members
                - Proposed Time
                    - Given by squad members
            - Filters:
                - Proposed Location
                    - Weighs locations with least amount of driving higher than ones with more.
                    - If there is a tie weighs locations closer to `target_loc` higher than those 
                      farther away.
                - Proposed Time
                    - Weights times based on how many conflicts there are with other peoples 
                      schedules
                    - People are given the times and then mark any of them if they don't work as a conflict
                    - The times with the least amount of conflicts are weighted the most
                    - If there is a tie then the times with the least amount of conflict for drivers 
                      are weighed higher.
        - `Movie`
            - Extends `Get Together`
            - Meta option: experience = standard | 3d | imax
            - Sources:
                - Proposed Location
                    - Use movie show times apis to get theatres showing movie near `target_loc`
                - Proposed Time
                    - Use movie show time apis to get show time options from locations
                - `data.experience`
                    - Data to compare against `meta.experience`
                    - Use movie show time apis to figure out if movie is Standard, 3D, or IMAX
            - Filters:
                - Proposed Location
                    - Compare `data.experience` of location against preferred `meta.experience`
                    - Weigh locations that provide the desired experience higher than those that don't
                - Proposed Time
                    - Compare `data.experience` of time against preferred `meta.experience`
                    - Weigh show times that provide the desired experience higher than those that don't
        - `Meal`
            - Extends `Get Together`
            - Meta option: type = fastfood | asian | drinks | mexican | ect.. (Any keywords to describe food types)
            - Sources
                - Proposed Location
                    - Use the Google Maps API to search for food with keywords
            - Filters
                - Proposed Time
                    - Weighs times where restaurants are at low peak and open higher than busy or closed.
    - `meta_option`: `string`
        - Depending on type this value is unmarshalled differently
        
### Event Proposal
A proposal for a place and time of the event. Event Proposals are generated by EventPlanners in the data source stage and 
then weighted in the filtering stage.

- `event_id`: `int`
- `weight`: `int`
    - Range of 0 to 100. Default to `0`.
- `loc_lat`: `long`
    - Proposed location latitude.
- `loc_long`: `long`
    - Proposed location longitude.
- `start_datetime`: `datetime`
    - Proposed datetime to start the event.
- `end_datetime`: `datetime`
    - Proposed datetime to end the event.

### Event RSVP
Records whether or not a user is going to an event.

- `event_id`: `int`
- `user_id`: `int`
- `attending`: `boolean`
- `active_car_id`: `int`
    - The id of the car that can be used by the user, `-1` if no car is available.

### Event Time Marker
Records a users ability to attend or not attend a certain proposed time 

- `event_id`: `int`
- `user_id`: `int`
- `event_proposal_id`: `int`
- `available`: `boolean`
- `active_car_id`: `int`
    - The id of the car that can be used by the user, `-1` if no car is available.
    
# Debt
Keeps track of debts to other users.

## Models
### Debt
Represents an entry in the overall "finacial logbook" of a user

- `payee_user_id`: `int`
    - Id of user that owes or is paying money.
    - Id of user who will be the source of the money.
- `receiver_user_id`: `int`
    - Id of user receiving money.
- `value`: `float`
    - If value is positive then payee user is paying a debt off to receiver
    - If value is negative then payee user owes receiver the value

# Polls
Polls present a set of options which users can vote for. Polls can be attached to squads or events

## Models
### Poll
Represents a set of questions users can interact with.

- `owner_type`: `string`
    - Either `squad` or `event` depending on if an event or squad owns the poll. 
    - The term "owns" is used to signify that the poll is displayed as part of a squad or as part of an event. The `creator_id` 
      field below specifies the user who created the poll.
- `owner_id`: `int`
    - Id of either squad or event depending on `owner_type`'s value.
- `creator_id`: `int`
    - Id of user who created poll
    
### Poll Option
Individual option for poll.

- `poll_id`: `int`
    - Id of poll option is for.
- `order`: `int`
    - Order in overall poll that option is in
- `name`: `string`
    - Name of the poll option
- `type`: `string`
    - Type of poll, put in place to support future more complicated polls
    - Currently only valid value is `text`
- `data`: `string`
    - Marshalled data for poll option, interpreted based on `type`.
    
### Poll Vote
Vote for poll option by user.

- `user_id`: `int`
- `poll_id`: `int`
- `poll_option_id`: `int`
    - Id of poll option they voted for
