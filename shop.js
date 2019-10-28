const APIGateway = 'https://p4b43mv7al.execute-api.us-west-2.amazonaws.com/dev'
const productID = 'prod_G3nbhaoJkINZ5v';

function fetchProduct() {
  return fetch(`${APIGateway}/products/${productID}`)
    .then((response) => {
      if (!response.ok) {
        throw 'failed to get product';
      }
      return response.json();
    })
    .then((json) => {
      console.log(json);
    })
    .catch((err) => console.log(err));
}

function fetchCoupons() {
  return fetch(`${APIGateway}/coupons`)
    .then((response) => {
      if (!response.ok) {
        throw 'failed to get coupons';
      }
      return response.json();
    })
    .then((json) => {
      console.log(json);
    })
    .catch((err) => console.log(err));
}

function responseHandler(response, url) {
  if (!response.ok) {
    throw `failed to GET ${url}`;
  }
  return response.json();
}

function formatCentPrice(cents) {
  const dollars = cents / 100;
  return dollars.toLocaleString('en-US', { style: 'currency', currency: 'USD' });
}

function buildProductGrid([product, { offers, coupons }]) {
  const {
    name,
    SKUList,
  } = product;

  // Tuna is redudant for now.
  const productName = name.replace(/Tuna /, '')

  // Only supports one shop keeping unit for now.
  // In the future there could be more flavors etc.
  const {
    price: productPrice,
    quantity: remainingQuantity,
  } = SKUList[0];

  const grid = document.getElementById('product-grid');
  grid.innerHTML = '';

  offers.forEach(({ quantity, couponID }) => {
    const title = document.createElement('strong');
    const price = document.createElement('em');
    const input = document.createElement('input');

    if (quantity == 1) title.innerHTML = `${quantity} ${productName}`;
    else title.innerHTML = `${quantity} ${productName}s`;

    // Stripe stores amounts as cent values.
    const grossAmount = productPrice * quantity;
    const centAmount = coupons[couponID] ? grossAmount - coupons[couponID] : grossAmount;
    price.innerHTML = formatCentPrice(centAmount);

    input.setAttribute('type', 'radio');
    input.setAttribute('name', 'product-select');
    input.setAttribute('value', quantity);

    // Fade the elements in as they are added to the DOM.
    [title, price, input].forEach((e) => {
      e.classList.add('fade-new-element');
      grid.append(e);
    });
  });
}

// Fetch the information for the product selection.
Promise.all([
  `${APIGateway}/products/${productID}`,
  `${APIGateway}/coupons`,
].map((url) => fetch(url)
  .then((response) => responseHandler(response, url))
  .catch((err) => console.log(err)))).then(buildProductGrid);

  /*
     <strong>1 Can</strong>
     <em> ...loading prices ... </em>
     <input type="radio" name="product-select" value="single">

     <strong>3 Cans</strong>
     <em> ...loading prices ... </em>
     <input type="radio" name="product-select" value="three">

     <strong>12 Cans</strong>
     <em> ...loading prices ... </em>
     <input type="radio" name="product-select" value="half">

     <strong>24 Cans</strong>
     <em> ...loading prices ... </em>
     <input type="radio" name="product-select" value="full">
   */
