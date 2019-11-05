<?php
$selectedLink = 'shop.html';
$jsHandler = 'confirm.js';
include('_header.php');
?>
<section>
  <noscript><h4 style="color:red">Enable scripts to use the shop</h4></noscript>
  <p>
    You will be emailed with the tracking number once the order is processed.
    Please review your order before submitting payment.
  </p>
  <div id="confirm-order">
    <em>loading order...</em>
    <div class="spinner"></div>
  </div>
  <form id="confirm-form">
    <div id="card-element"><!-- Stripe iframe --></div>
    <div id="submit-row" class="info-row">
      <button type="submit">Place Order</button>
      <span id="form-errors" for="place-order" role="alert"></label>
    </div>
  </form>
</section>
<?php include('_footer.php') ?>
