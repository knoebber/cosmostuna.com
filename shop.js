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
    .catch((err) => console.log(err));
}

function responseHandler(response, url) {
  if (!response.ok) {
    throw `failed request to ${url}`;
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
    id: SKUID,
    price: productPrice,
    quantity: remainingQuantity, // TODO use this to remove unavailable products.
  } = SKUList[0];

  const productGrid = document.getElementById('product-grid');
  productGrid.innerHTML = '';
  productGrid.setAttribute('data-sku-id', SKUID)

  offers.forEach(({ quantity, couponID }) => {
    const title = document.createElement('label');
    const price = document.createElement('em');
    const input = document.createElement('input');

    if (quantity == 1) {
      title.innerHTML = `${quantity} ${productName}`;
      // Set buying one can to the default.
      input.setAttribute('checked', true);
    } else {
      title.innerHTML = `${quantity} ${productName}s`;
    }

    const target = `input-quantity-${quantity}`;
    title.setAttribute('for', target)

    // Stripe stores amounts as cent values.
    const grossAmount = productPrice * quantity;
    const centAmount = coupons[couponID] ? grossAmount - coupons[couponID] : grossAmount;
    price.innerHTML = formatCentPrice(centAmount);

    input.setAttribute('type', 'radio');
    input.setAttribute('name', 'product-select');
    input.setAttribute('id', target);
    input.setAttribute('value', quantity);

    // Fade the elements in as they are added to the DOM.
    [title, price, input].forEach((e) => {
      e.classList.add('fade-new-element');
      productGrid.append(e);
    });
  });
}

function placeOrder(token) {
  const productGrid = document.getElementById('product-grid');
  const SKUID = productGrid.getAttribute('data-sku-id');
  const quantity = Array.from(productGrid.children).filter(e => e.checked)[0].value;
  
  const formValues = [
    'name',
    'address',
    'city',
    'zip',
    'state',
    'email',
  ].reduce((result, curr) => ({
    [curr]: document.getElementById(curr).value,
    ...result
  }), {});

  const orderBody = {
    email: formValues.email,
    shipping: {
      name: formValues.name,
      address: {
        city: formValues.city,
        line1: formValues.address,
        postal_code: formValues.zip,
        state: formValues.state,
      },
    },
    items: [{
      quantity: parseInt(quantity, 10),
      parent: SKUID,
    }],
  };

  const url = `${APIGateway}/orders`;

  fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(orderBody),
  })
    .then((response) => responseHandler(response, url))
    .then(({ target, message, orderID }) => {
      if (orderID) {
        window.location.href = `/confirm.html?order=${orderID}`
      } else {
        const formErrors = document.getElementById('form-errors');
        formErrors.innerHTML = message;
        // TODO highlight target.
      }
    });
}

function loadShop() {
  // Fetch the information for the product selection.
  Promise.all([
    `${APIGateway}/products/${productID}`,
    `${APIGateway}/coupons`,
  ].map((url) => fetch(url)
    .then((response) => responseHandler(response, url))
    .catch((err) => console.log(err)))).then(buildProductGrid);

  const form = document.getElementById('shop-form');

  form.addEventListener('submit', (event) => {
    event.preventDefault();
    placeOrder()
  });
}

window.onload = loadShop;
