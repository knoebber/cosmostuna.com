const APIGateway = 'https://p4b43mv7al.execute-api.us-west-2.amazonaws.com/dev'
const productID = 'prod_G3nbhaoJkINZ5v';
const stripePublicKey = 'pk_test_qUdrpKjmC5gZ7jcuuHeRb8Au006WnfLwAt';

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
    const title = document.createElement('label');
    const price = document.createElement('em');
    const input = document.createElement('input');

    if (quantity == 1) title.innerHTML = `${quantity} ${productName}`;
    else title.innerHTML = `${quantity} ${productName}s`;

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
      grid.append(e);
    });
  });
}

function placeOrder(token) {
  const address = [
    'name',
    'address',
    'city',
    'zip',
    'state',
  ].reduce((result, curr) => ({
    [curr]: document.getElementById(curr).value,
    ...result
  }), {});

  console.log(address);

  console.log(token);
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
  const submitButton = document.getElementById('place-order');
  const formErrors = document.getElementById('form-errors');

  // Initialize Stripe.
  const stripe = Stripe(stripePublicKey);
  const elements = stripe.elements();

  const style = {
    base: {
      fontSize: '20px',
      color: '#32325d',
    },
  };

  const card = elements.create('card', { style });
  card.mount('#card-element');

  card.addEventListener('change', ({ error }) => {
    if (error) {
      submitButton.setAttribute('disabled', true);
      formErrors.textContent = error.message;
    } else {
      submitButton.setAttribute('disabled', false);
      formErrors.textContent = '';
    }
  });

  form.addEventListener('submit', async (event) => {
    event.preventDefault();

    const { token, error } = await stripe.createToken(card);

    if (error) {
      formErrors.textContent = error.message;
    } else {
      placeOrder(token);
    }
  });
}

window.onload = loadShop;
