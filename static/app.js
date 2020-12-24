const Controller = {
  query: "",
  page: 1,
  load: () => {
    const form = document.getElementById("form");
    const prev = document.querySelector(".prev");
    const next = document.querySelector(".next");
    form.addEventListener("submit", Controller.search);
    prev.addEventListener("click", () => Controller.makeRequest(Controller.query, Controller.page - 1));
    next.addEventListener("click", () => Controller.makeRequest(Controller.query, Controller.page + 1));
  },

  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    Controller.makeRequest(data.query)
  },

  makeRequest: (query, page = 1) => {
    const response = fetch(`/search?q=${query}&p=${page}`).then((response) => {
      response.json().then((results) => {
        Controller.page = page;
        Controller.query = query;
        Controller.updateTable(results, query, page);
      })
    })
    .catch((err) => {
      console.log(err)
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const prev = document.querySelector(".prev");
    const next = document.querySelector(".next");

    if (Controller.page == 1) {
      prev.style.display = "none";
    } else if (Controller.page > 1) {
      prev.style.display = "block";
    }

    if (results.length < 20) {
      next.style.display = "none";
    } else {
      next.style.display = "block";
    }

    const rows = [];
    const regex = RegExp(`(${Controller.query})`, 'gmi');

    for (let result of results) {
      rows.push(
        `<tr>
          <td>
            <pre>${result.trim().replace(regex, `<mark>$1</mark>`)}</pre>
          </td>
        <tr/>`
      );
    }

    table.innerHTML = rows.join("");
  },
};

Controller.load();