create or replace function dbo.try_lock_seats(
    p_flight_id uuid,
    p_quantity int
)
returns table (
    seat_lock_id int
)
language plpgsql
as $$
declare
    v_max_available_seats int;
    v_existing_lock_count int;
begin

    if p_quantity <= 0 then
        raise exception 'Quantity must be greater than 0';
    end if;

    /*
        Lock the flight row so all seat allocations for the
        same flight are serialized safely.
    */
    select max_available_seats
    into v_max_available_seats
    from dbo.flight
    where id = p_flight_id
    for update;

    if not found then
        raise exception 'Flight not found: %', p_flight_id;
    end if;

    /*
        Count existing seat locks
    */
    select count(*)
    into v_existing_lock_count
    from dbo.seat_lock
    where flight_id = p_flight_id;

    /*
        Ensure enough remaining capacity exists
    */
    if v_existing_lock_count + p_quantity > v_max_available_seats then
        raise exception
            'Insufficient seats available. Requested %, available %',
            p_quantity,
            v_max_available_seats - v_existing_lock_count;
    end if;

    /*
        Insert N seat locks and return their IDs
    */
    return query
    insert into dbo.seat_lock (flight_id)
    select p_flight_id
    from generate_series(1, p_quantity)
    returning id;

end;
$$;
