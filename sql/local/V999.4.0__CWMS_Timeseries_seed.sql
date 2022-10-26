
DO $$

DECLARE 
    lrn_cwms_timeseries_ds_id uuid;
    mvp_cwms_timeseries_ds_id uuid;
    nwo_cwms_timeseries_ds_id uuid;
    nwp_cwms_timeseries_ds_id uuid;
    nws_cwms_timeseries_ds_id uuid;

BEGIN

    SELECT d.id into lrn_cwms_timeseries_ds_id
    FROM datasource d 
    JOIN datasource_type dt ON dt.id = d.datasource_type_id 
    JOIN provider p ON p.id = d.provider_id 
    WHERE dt.slug = 'cwms-timeseries' AND p.slug = 'lrn';

    SELECT d.id into mvp_cwms_timeseries_ds_id
    FROM datasource d 
    JOIN datasource_type dt ON dt.id = d.datasource_type_id 
    JOIN provider p ON p.id = d.provider_id 
    WHERE dt.slug = 'cwms-timeseries' AND p.slug = 'mvp';

    SELECT d.id into nwo_cwms_timeseries_ds_id
    FROM datasource d 
    JOIN datasource_type dt ON dt.id = d.datasource_type_id 
    JOIN provider p ON p.id = d.provider_id 
    WHERE dt.slug = 'cwms-timeseries' AND p.slug = 'nwo';

    SELECT d.id into nwp_cwms_timeseries_ds_id
    FROM datasource d 
    JOIN datasource_type dt ON dt.id = d.datasource_type_id 
    JOIN provider p ON p.id = d.provider_id 
    WHERE dt.slug = 'cwms-timeseries' AND p.slug = 'nwp';

    SELECT d.id into nws_cwms_timeseries_ds_id
    FROM datasource d 
    JOIN datasource_type dt ON dt.id = d.datasource_type_id 
    JOIN provider p ON p.id = d.provider_id 
    WHERE dt.slug = 'cwms-timeseries' AND p.slug = 'nws';




    INSERT into timeseries(datasource_id, datasource_key, latest_time, latest_value) VALUES
    -- LRN cwms-timeseries
    (lrn_cwms_timeseries_ds_id, 'BAHK2-BARKLEY.Elev-Pool.Inst.15Minutes.0.dcp-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BAHK2-BARKLEY.Elev-Pool.Inst.1Hour.0.man-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BAHK2-BARKLEY.Elev-Tail.Inst.1Hour.0.man-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BAHK2-BARKLEY.Flow.Ave.1Hour.1Hour.man-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BAHK2-BARKLEY.Flow-In.Ave.1Hour.1Hour.man-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BAHK2-BARKLEY.Precip-Inc.Total.15Minutes.15Minutes.dcp-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BAHK2-BARKLEY.Precip-Inc.Total.~1Day.1Day.dcp-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BAHK2-BARKLEY.Precip-Inc.Total.1Hour.1Hour.dcp-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BAHK2-BARKLEY.Precip-Inc.Total.~6Hours.6Hours.dcp-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BAHK2-BARKLEY.Stor.Inst.1Hour.0.man-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BARK2-BARKLEY.Elev-Tail.Inst.1Hour.0.man-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BARK2-BARKLEY.Energy.Total.1Hour.1Hour.man-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BARK2-BARKLEY.Flow.Ave.1Hour.1Hour.man-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BARK2-BARKLEY.Flow.Ave.6Hours.6Hours.tva-qpf', null, null),
    (lrn_cwms_timeseries_ds_id, 'BARK2-BARKLEY.Flow-Spillway.Ave.1Hour.1Hour.man-rev', null, null),
    (lrn_cwms_timeseries_ds_id, 'BARK2-BARKLEY.Flow-Turbine.Ave.1Hour.1Hour.man-rev', null, null),
    -- MVP cwms-timeseries
    (mvp_cwms_timeseries_ds_id, 'Baldhill_Dam.Elev.Inst.15Minutes.0.rev-NGVD29', null, null),
    (mvp_cwms_timeseries_ds_id, 'Baldhill_Dam.Flow-In.Ave.15Minutes.1Day.comp', null, null),
    (mvp_cwms_timeseries_ds_id, 'Baldhill_Dam.Flow-Out.Inst.15Minutes.0.rev', null, null),
    (mvp_cwms_timeseries_ds_id, 'Baldhill_Dam.Precip-cum.Inst.15Minutes.0.rev', null, null),
    (mvp_cwms_timeseries_ds_id, 'Baldhill_Dam.Stage.Inst.15Minutes.0.rev', null, null),
    (mvp_cwms_timeseries_ds_id, 'Baldhill_Dam.Stor.Inst.15Minutes.0.comp', null, null),
    (mvp_cwms_timeseries_ds_id, 'Baldhill_Dam-Tailwater.Elev.Inst.15Minutes.0.rev-NAVD88', null, null),
    (mvp_cwms_timeseries_ds_id, 'Baldhill_Dam-Tailwater.Elev.Inst.15Minutes.0.rev-NGVD29', null, null),
    (mvp_cwms_timeseries_ds_id, 'Baldhill_Dam-Tailwater.Stage.Inst.15Minutes.0.rev', null, null),
    (mvp_cwms_timeseries_ds_id, 'Baldhill_Dam-Tailwater.Temp-Water.Inst.15Minutes.0.CEMVP-GOES-Raw', null, null),
    (mvp_cwms_timeseries_ds_id, 'Baldhill_Dam.Temp-Air.Inst.15Minutes.0.CEMVP-GOES-Raw', null, null),
    -- NWO cwms-timeseries
    (nwo_cwms_timeseries_ds_id, 'BECR-Bear_Creek_Dam-Bear.Elev.Inst.1Hour.0.Reporting', null, null),
    (nwo_cwms_timeseries_ds_id, 'BECR-Bear_Creek_Dam-Bear.Stor.Inst.1Hour.0.Reporting', null, null),
    (nwo_cwms_timeseries_ds_id, 'BECR.Elev.Inst.1Hour.0.Best-NWO', null, null),
    (nwo_cwms_timeseries_ds_id, 'BECR.Flow-In.Ave.1Hour.6Hours.Best-NWO', null, null),
    (nwo_cwms_timeseries_ds_id, 'BECR.Flow-In.Inst.1Hour.0.Best-NWO', null, null),
    (nwo_cwms_timeseries_ds_id, 'BECR.Flow-Out.Inst.1Hour.0.Best-NWO', null, null),
    (nwo_cwms_timeseries_ds_id, 'BECR.Stor.Inst.1Hour.0.Best-NWO', null, null),
    (nwo_cwms_timeseries_ds_id, 'BECR-Surface.Temp-Water.Inst.1Hour.0.Rev-NWO-Evap', null, null),
    -- NWP cwms-timeseries
    (nwp_cwms_timeseries_ds_id, 'APP.Elev-Forebay.Inst.0.0.Best', null, null),
    (nwp_cwms_timeseries_ds_id, 'APP.Elev-Forebay.Inst.0.0.MIXED-REV', null, null),
    (nwp_cwms_timeseries_ds_id, 'APP.Flow-In.Ave.15Minutes.6Hours.Best', null, null),
    (nwp_cwms_timeseries_ds_id, 'APP.Flow-In.Ave.15Minutes.6Hours.MIXED-COMPUTED-REV', null, null),
    (nwp_cwms_timeseries_ds_id, 'APP.Flow-Out.Inst.0.0.Best', null, null),
    (nwp_cwms_timeseries_ds_id, 'APP.Flow-Out.Inst.0.0.MIXED-COMPUTED-REV', null, null),
    (nwp_cwms_timeseries_ds_id, 'APP.Stor.Inst.0.0.Best', null, null),
    (nwp_cwms_timeseries_ds_id, 'APP.Stor.Inst.0.0.MIXED-COMPUTED-REV', null, null),
    (nwp_cwms_timeseries_ds_id, 'BCL.Elev-Forebay.Inst.0.0.Best', null, null),
    (nwp_cwms_timeseries_ds_id, 'BCL.Elev-Forebay.Inst.0.0.MIXED-REV', null, null),
    (nwp_cwms_timeseries_ds_id, 'BCL.Flow-In.Ave.15Minutes.6Hours.Best', null, null),
    (nwp_cwms_timeseries_ds_id, 'BCL.Flow-In.Ave.15Minutes.6Hours.MIXED-COMPUTED-REV', null, null),
    (nwp_cwms_timeseries_ds_id, 'BCL.Flow-Out.Inst.0.0.Best', null, null),
    (nwp_cwms_timeseries_ds_id, 'BCL.Flow-Out.Inst.0.0.MIXED-COMPUTED-REV', null, null),
    -- NWS cwms-timeseries
    (nws_cwms_timeseries_ds_id, 'ALF.Elev-Forebay.Inst.1Hour.0.CBT-RAW', null, null),
    (nws_cwms_timeseries_ds_id, 'ALF.Elev-Forebay.Inst.1Hour.0.CBT-REV', null, null),
    (nws_cwms_timeseries_ds_id, 'ALF.Flow-Out.Ave.1Hour.1Hour.CBT-RAW', null, null),
    (nws_cwms_timeseries_ds_id, 'ALF.Flow-Out.Ave.1Hour.1Hour.CBT-REV', null, null),
    (nws_cwms_timeseries_ds_id, 'ALF.Stor-Lake.Inst.1Hour.0.NWSRADIO-COMPUTED-RAW', null, null),
    (nws_cwms_timeseries_ds_id, 'ALF.Stor.Inst.~6Hours.0.MODEL-STP-FCST', null, null),
    (nws_cwms_timeseries_ds_id, 'HAH.Elev-Forebay.Inst.1Hour.0.NWSRADIO-REV', null, null),
    (nws_cwms_timeseries_ds_id, 'HAH.Flow-In.Inst.1Hour.0.NWSRADIO-COMPUTED-REV', null, null),
    (nws_cwms_timeseries_ds_id, 'HAH.Precip-Cum.Inst.1Hour.0.NWSRADIO-REV', null, null),
    (nws_cwms_timeseries_ds_id, 'HAH.Precip-Inc.Total.1Hour.1Hour.NWSRADIO-REV', null, null),
    (nws_cwms_timeseries_ds_id, 'HAH.Stor.Inst.1Hour.0.NWSRADIO-REV', null, null),
    (nws_cwms_timeseries_ds_id, 'HAH.Temp-Air.Inst.1Hour.0.NWSRADIO-REV', null, null);

END$$;
