const APIGateway = 'https://p4b43mv7al.execute-api.us-west-2.amazonaws.com/dev'
const stripePublicKey = 'pk_test_qUdrpKjmC5gZ7jcuuHeRb8Au006WnfLwAt';

function loadConfirm() {
  const form = document.getElementById('confirm-form');
  const formErrors = document.getElementById('form-errors');
  const orderID = new URLSearchParams(window.location.search).get('order');

  // Load the order.
  fetch(`${APIGateway}/orders/${orderID}`)
    .then((response) => {
      if (!response.ok) {
        throw 'failed to get order';
      }
      return response.json();
    })
    .then(({
      email,
      amount,
      items,
      shipping: {
        name,
        address: {
          city,
          country,
          line1,
          postal_code: zip,
          state,
        },
      },
    }) => {
      const confirmOrder = document.getElementById('confirm-order');
      confirmOrder.innerHTML = '';

      const { description, quantity } = items[0];
      [
        {
          name: 'Name',
          value: name,
        },
        {
          name: 'Address',
          value: line1,
        },
        {
          name: 'City',
          value: `${city}, ${state} ${zip}`,
        },
        {
          name: 'Email',
          value: email,
        },
        {
          name: 'Description',
          value: `${quantity} ${quantity > 1 ? 'cans': 'can'} tuna`},
        {
          name: 'Amount Due',
          value: (amount / 100).toLocaleString('en-US', { style: 'currency', currency: 'USD' }),
        },
      ].forEach(({ name, value }) => {
        const infoRow = document.createElement('div');
        infoRow.classList.add('info-row', 'fade-new-element');

        const title = document.createElement('strong');
        title.innerHTML = name;

        const text = document.createElement('span');
        text.innerHTML = value;

        infoRow.append(title, text);
        confirmOrder.append(infoRow);
      });
    })
    .catch((err) => {
      console.log(err)
      formErrors.textContent = 'Unable to retrieve order.';
      document.getElementsByTagName('button')[0].remove()
    });

  // Initialize Stripe.
  const stripe = Stripe(stripePublicKey);
  const elements = stripe.elements();

  const style = {
    base: {
      fontFamily: 'serif',
      fontSize: '19px',
      color: '#232f3e',
    },
  };

  const card = elements.create('card', { style });
  card.mount('#card-element');

  card.addEventListener('change', ({ error }) => {
    if (error) {
      formErrors.textContent = error.message;
    } else {
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

window.onload = loadConfirm;
