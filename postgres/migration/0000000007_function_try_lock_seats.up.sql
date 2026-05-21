create or replace function dbo.try_lock_seats(
    p_flight_id uuid,
    p_quantity int
)
returns table (
    flight_valid boolean,
    seats_available boolean,
    seat_lock_ids int[]
)
language plpgsql
as $$
declare
    v_max_available_seats int;
    v_existing_lock_count int;
    v_inserted_ids int[];
begin
    -- Protect against invalid inputs
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

    -- Scenario 1: Flight does not exist
    if not found then
        return query select false, false, null::int[];
        return;
    end if;

    /*
        Count existing seat locks
    */
    select count(*)
    into v_existing_lock_count
    from dbo.seat_lock
    where flight_id = p_flight_id;

    -- Scenario 2: Flight exists, but not enough capacity
    if v_existing_lock_count + p_quantity > v_max_available_seats then
        return query select true, false, null::int[];
        return;
    end if;

    /*
        Scenario 3: Success. Insert N seat locks, collect IDs into an array,
        and return the final payload row.
    */
    with inserted_rows as (
        insert into dbo.seat_lock (flight_id)
        select p_flight_id
        from generate_series(1, p_quantity)
        returning id
    )
    select array_agg(id) into v_inserted_ids from inserted_rows;

    return query select true, true, v_inserted_ids;
end;
$$;

