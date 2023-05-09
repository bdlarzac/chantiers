/** 
    https://stackoverflow.com/questions/14267781/sorting-html-table-with-javascript
    Solution of Cesar Morillas
    
    Also used https://stackoverflow.com/questions/18452743/localcompare-for-integers
    Answer of pomo
    
    Usage : include in html page <script src="table_sort.js"></script> AFTER the tables (or use defer).
    The call to table_sort() (see end of current file) must be done after the tables are displayed.
    
**/
function table_sort() {
  const styleSheet = document.createElement('style')
  styleSheet.innerHTML = `
        .order{
            color:blue;
        }
        .order-inactive span {
            visibility:hidden;
        }
        .order-inactive:hover span {
            visibility:visible;
        }
        .order-active span {
            visibility: visible;
        }
        .order{
            cursor: pointer;
        }
    `
  document.head.appendChild(styleSheet);

  document.querySelectorAll('th.order').forEach(th_elem => {
    let asc = true;
    const span_elem = document.createElement('span');
    span_elem.style = "font-size:0.8rem; margin-left:0.5rem";
    span_elem.innerHTML = "▼";
    th_elem.appendChild(span_elem);
    th_elem.classList.add('order-inactive');

    const index = Array.from(th_elem.parentNode.children).indexOf(th_elem);
    
    th_elem.addEventListener('click', (e) => {
      document.querySelectorAll('th.order').forEach(elem => {
        elem.classList.remove('order-active');
        elem.classList.add('order-inactive');
      })
      th_elem.classList.remove('order-inactive');
      th_elem.classList.add('order-active');

      if (!asc) {
        th_elem.querySelector('span').innerHTML = '▲';
      } else {
        th_elem.querySelector('span').innerHTML = '▼';
      }
      const arr = Array.from(th_elem.closest("table").querySelectorAll('tbody tr'));
      arr.sort((a, b) => {
        // const a_val = a.children[index].innerText;
        // const b_val = b.children[index].innerText;
        // note Thierry : changé innerText en innerHTML de manière à pouvoir trier par date
        // avec un hack du style : <span data-date="{{.DateActivite}}">{{.DateActivite | dateFr}}</span>
        const a_val = a.children[index].innerHTML;
        const b_val = b.children[index].innerHTML;
        return (asc) ? a_val.localeCompare(b_val, undefined, {'numeric': true}) : b_val.localeCompare(a_val, undefined, {'numeric': true});
      });
      arr.forEach(elem => {
        th_elem.closest("table").querySelector("tbody").appendChild(elem);
      });
      asc = !asc;
    })
  })
}

table_sort();
