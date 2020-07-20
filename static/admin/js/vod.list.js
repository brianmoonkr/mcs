$(document).ready(function () {
  const $vodListTable = $('#vod-list-table');
  const $vodForm = $('#vod-form');
  const $status = $('#status');
  const $dateType = $('#date-type');
  const $startDate = $('#start-date');
  const $endDate = $('#end-date');
  const $searchTextType = $('#search-text-type');
  const $searchText = $('#search-text');
  const $searchBtn = $("#search-btn");
  const $searchDefaultBtn = $("#search-default-btn");

  var SearchData = {
    status: "",
    dateType: "",
    startDate: null,
    endDate: null,
    searchTextType: "",
    searchText: ""
  };

  // Datepicker
  $startDate.datepicker({
    uiLibrary: 'bootstrap4',
    format: 'yyyy-mm-dd'
  });
  $endDate.datepicker({
    uiLibrary: 'bootstrap4',
    format: 'yyyy-mm-dd'
  });

  const vodListTable = $vodListTable.DataTable({
    responsive: true,
    paging: true,
    lengthChange: true,
    searching: false,
    ordering: false,
    info: true,
    autoWidth: true,
    processing: true,
    serverSide: true,
    language: {
      infoFiltered: "",
      info: "Showing _START_ to _END_ of _TOTAL_ entries"
    },
    sPaginationType: "full_numbers",
    ajax: {
      url: "/api/v1/admin/vod/list",
      contentType: 'application/json',
      type: "POST",
      dataSrc: "admin_vod_list",
      data: function (data) {
        data.status = SearchData.status;
        data.start_date = SearchData.startDate;
        data.end_date = SearchData.endDate;
        data.date_type = SearchData.dateType;
        data.search_text_type = SearchData.searchTextType;
        data.search_text = SearchData.searchText;
        return JSON.stringify(data);
      }
    },
    columns: [
      { data: 'row_num' },
      { data: 'title' },
      { data: 'nick_name' },
      { data: 'view_cnt' },
      { data: 'like_cnt' },
      { data: 'duration' },
      { data: 'created_at' },
    ],
    columnDefs: [
      {
        render: function (data, type, row) {
          return `<a href="/admin/vod/detail/${row.vod_seq}">${data}</a>`;
        },
        targets: [1]
      },
      {
        render: function (data, type, row) {
          //return moment('2000-01-01 00:00:00').add(moment.duration(Number(data))).format('HH:mm:ss');
          return moment.utc(Number(data)).format('HH:mm:ss');
        },
        targets: [5]
      },
      {
        render: function (data, type, row) {
          if (data.includes("0001-")) {
            return "-";
          } else {
            return moment(data).format('YYYY-MM-DD HH:mm:ss');
          }
        },
        targets: [6]
      }
    ]
  });


});