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
    <hr>
    <div id="disclaimer-row" style="display:none;">
      <p>By paying <span id="total-amount"></span> you agree to our <a target="_blank" href="/terms.html">terms.</a></p>
    </div>
    <div id="card-element"><!-- Stripe iframe --></div>
    <div id="submit-row" class="info-row">
      <button type="submit">Place Order</button>
      <span id="form-errors" for="place-order" role="alert"></label>
    </div>
  </form>
</section>
