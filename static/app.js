const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    fetch(`/search?q=${data.query}&page=${data.page}&limit=${data.limit}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    let index = 0;
    table.innerHTML = "";
    Controller.addCell(table, index, `Total results: ${results.total_quantity} <br> Page results: ${results.page_quantity}`);
    index++;
    for (let result of results.search_results) {
      Controller.addCell(table, index, `Lines: ${result.from_line} to ${result.to_line}`);
      index++;
      let text = "";
      for (let line of result.lines_text) {
        text += `${line}<br>`;
      }
      Controller.addCell(table, index, text);
      index++;
    }
  },

  addCell: (table, index, html) => {
    let row_line = table.insertRow(index);
    let cell_line = row_line.insertCell(0);
    cell_line.innerHTML = html;
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
