const productID = 'prod_G3nbhaoJkINZ5v';

function fetchProduct() {
  return fetch(`${apiGateway}/products/${productID}`)
    .then((response) => {
      if (!response.ok) {
        throw 'failed to get product';
      }
      return response.json();
    })
    .catch((err) => console.log(err));
}

function fetchCoupons() {
  return fetch(`${apiGateway}/coupons`)
    .then((response) => {
      if (!response.ok) {
        throw 'failed to get coupons';
      }
      return response.json();
    })
    .catch((err) => console.log(err));
}

// The arguments are destructuring the list of promise.all() results.
function buildProductGrid([productData, couponData]) {
  const productGrid = document.getElementById('product-grid');

  if (!productData || !couponData) {
    productGrid.innerHTML = 'Failed to fetch products.';
    document.getElementById('shop-form').innerHTML = '';
    return;
  } else {
    productGrid.innerHTML = '';
  }

  const { name, SKUList } = productData;
  const { offers, coupons } = couponData;

  // Tuna is redudant for now.
  const productName = name.replace(/Tuna /, '')

  // Only supports one shop keeping unit for now.
  // In the future there could be more flavors etc.
  const {
    id: SKUID,
    price: productPrice,
    quantity: productsLeft,
  } = SKUList[0];

  productGrid.setAttribute('data-sku-id', SKUID)
  if (productsLeft < 1){
    productGrid.innerHTML = 'Out of stock.';
    document.getElementById('order-grid').remove()
    document.getElementById('submit-row').remove()
    return;
  }

  offers.forEach(({ quantity, couponID }) => {
    if (productsLeft < quantity) {
      return;
    }

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
        postalCode: formValues.zip,
        state: formValues.state,
      },
    },
    items: [{
      quantity: parseInt(quantity, 10),
      parent: SKUID,
    }],
  };

  const url = `${apiGateway}/orders`;

  setDisabled('button', true)
  fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(orderBody),
  })
    .then((response) => responseHandler(response, url))
    .then(({ message, target, orderID }) => {
      if (orderID) {
        window.location.href = prod? `/confirm.html?order=${orderID}` : `/dev-confirm.html?order=${orderID}`;
      } else {
        formError(message, target);
      }
    })
    .finally(() => {
      setDisabled('button', false);
    });
}

function loadShop() {
  // Fetch the information for the product selection.
  setDisabled('button', true);
  Promise.all([
    `${apiGateway}/products/${productID}`,
    `${apiGateway}/coupons`,
  ].map((url) => fetch(url)
    .then((response) => responseHandler(response, url))
    .catch((err) => console.log(err))))
         .then(buildProductGrid)
         .finally(() => { setDisabled('button', false); });

  const form = document.getElementById('shop-form');

  form.addEventListener('submit', (event) => {
    event.preventDefault();
    placeOrder();
  });
}

window.onload = loadShop;
