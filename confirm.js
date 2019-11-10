const stripePublicKey = 'pk_test_qUdrpKjmC5gZ7jcuuHeRb8Au006WnfLwAt';

function submitPayment({ id: token }, orderID) {
  setDisabled('button', true);

  fetch(`${APIGateway}/orders`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      token,
      orderID,
    }),
  })
    .then((response) => {
      if (response.ok) {
        window.location.href = '/thank-you.html';
      } else {
        throw 'failed to submit order payment';
      }
    })
    .catch((err) => {
      console.log(err);
      formError('Unable to process payment')
    })
    .finally(() => {
      setDisabled('button', false);
    });
}

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
      status,
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
          value: formatCentPrice(amount),
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

      if (status !== 'created') {
        // TODO type switch on status, include link in email.
        formError('This order is already paid.');
        document.getElementsByTagName('button')[0].remove()
        return;
      }
    })
    .catch((err) => {
      console.log(err)
      formError('Unable to retrieve order');
      document.getElementsByTagName('button')[0].remove();
    });

  // Initialize Stripe.
  const stripe = Stripe(stripePublicKey);
  const elements = stripe.elements();

  const style = {
    base: {
      fontFamily: 'serif',
      fontSize: '20px',
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

    const { token, error } = await stripe.createToken(card, orderID);

    if (error) {
      formErrors.textContent = error.message;
    } else {
      submitPayment(token, orderID);
    }
  });
}

window.onload = loadConfirm;
