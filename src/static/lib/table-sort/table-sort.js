/** 
    https://stackoverflow.com/questions/14267781/sorting-html-table-with-javascript
    Solution of Cesar Morillas
    
    Usage : include in html page <script src="table_sort.js"></script> after the tables (or use defer).
    But the call to table_sort() must be done after.
    
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
        const a_val = a.children[index].innerText;
        const b_val = b.children[index].innerText;
        return (asc) ? a_val.localeCompare(b_val) : b_val.localeCompare(a_val);
      });
      arr.forEach(elem => {
        th_elem.closest("table").querySelector("tbody").appendChild(elem);
      });
      asc = !asc;
    })
  })
}

table_sort();
