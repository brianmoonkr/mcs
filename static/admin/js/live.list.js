$(document).ready(function () {
  const $liveListTable = $('#live-list-table');
  const $liveForm = $('#live-form');
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

  const liveListTable = $liveListTable.DataTable({
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
      url: "/api/v1/admin/live/list",
      contentType: 'application/json',
      type: "POST",
      dataSrc: "admin_live_list",
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
      { data: 'status' },
      { data: 'live_cnt' },
      { data: 'view_cnt' },
      { data: 'like_cnt' },
      { data: 'start_dt' },
      { data: 'end_dt' },
    ],
    columnDefs: [
      {
        render: function (data, type, row) {
          return `<a href="/admin/live/detail/${row.live_seq}">${data}</a>`;
        },
        targets: [1]
      },
      {
        render: function (data, type, row) {
          if (data == "1001") {
            return "등록";
          } else if (data == "1002") {
            return "대기";
          } else if (data == "1003") {
            return "방송중";
          } else if (data == "1004") {
            return "방송종료";
          } else if (data == "1005") {
            return "비정상종료";
          }
        },
        targets: [3]
      },
      {
        render: function (data, type, row) {
          if (data.includes("0001-")) {
            return "-";
          } else {
            return moment(data).format('YYYY-MM-DD HH:mm:ss');
          }
        },
        targets: [7, 8]
      }
    ]
  });


});