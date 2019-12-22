<section>
  <?php if(!$prod) :?>
    <div style="text-align: center;"><strong>[TEST MODE]</strong></div>
  <?php endif; ?>
  <noscript><h4 style="color:red">Enable scripts to use the shop</h4></noscript>
  <div id="confirm-action"></div>
  <hr>
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
