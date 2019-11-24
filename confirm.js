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

  const confirmAction = document.getElementById('confirm-action');
  const payButton = document.querySelector('#submit-row button');
  const cardForm = document.getElementById('card-element');
  const removePayment = ()=> {
    payButton.remove();
    cardForm.remove();
  };

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
        tracking_number: tracking,
        address: {
          city,
          country,
          line1,
          postal_code: zip,
          state,
        },
      },
    }) => {
      let paymentInfo;
      let trackingLink;
      const confirmP = document.createElement('p');
      if (status === 'created') {
        confirmP.textContent = 'Please review your order.';
        paymentInfo = 'Amount Due';
      } else if (status === 'paid') {
        confirmP.textContent = 'Your order is paid and pending shipping.';
        removePayment();
        paymentInfo = 'Amount Paid';
      } else if (status === 'fulfilled') {
        confirmP.textContent = 'Your order has been shipped.';
        if (tracking) {
          trackingLink = document.createElement('a');
          trackingLink.setAttribute('href', `https://tools.usps.com/go/TrackConfirmAction?tLabels=${tracking}`);
          trackingLink.setAttribute('target', '_blank');
          trackingLink.setAttribute('style', 'color:black;');
          trackingLink.textContent = `${tracking} (USPS)`;
        }
        paymentInfo = 'Amount Paid';
        removePayment();
      } else if(status === 'cancelled') {
        paymentInfo = 'Amount Refunded';
        confirmP.textContent = 'Your order is canceled.';
        removePayment();
      }
      confirmAction.appendChild(confirmP);
      if (trackingLink) {
        confirmAction.appendChild(trackingLink);
      }
        
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
          name: paymentInfo,
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

    })
    .catch((err) => {
      console.log(err)
      formError('Unable to retrieve order');
      removePayment();
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
