<?php
$apiGateway = 'https://p4b43mv7al.execute-api.us-west-2.amazonaws.com';
$stripeTestPublicKey = 'pk_test_qUdrpKjmC5gZ7jcuuHeRb8Au006WnfLwAt';
$stripeProdPublicKey = 'pk_live_yQTbrZw5CPfFetLlebxnTll7008omolXmg';

function buildLinks() {
  global $selectedLink;

  $links = array(
    '' => 'HOME',
    'shop.html' => 'BUY NOW',
    'gallery.html' => 'GALLERY',
    'about.html' => 'ABOUT',
  );

  $i = 0;
  foreach ($links as $path => $name) {
    if ($i++ === 2){
      echo '<img src="images/alpha-small-logo.png" alt="Logo">';
    }

    if ($selectedLink === $path){
      echo "<a href=\"/$path\" id=\"link-selected\">$name</a>";
    } else {
      echo "<a href=\"/$path\">$name</a>";
    }
    echo "\n";
  }
}
?>

<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
 "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" lang="en" xml:lang="en">
  <head>
    <meta http-equiv="Content-Type" content="text/html;charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <link rel="stylesheet" type="text/css" href="style.css" />
    <!-- Include Stripe script on all pages for fraud protection. -->
    <!-- https://stripe.com/docs/web/setup -->
    <script type="text/javascript" src="https://js.stripe.com/v3/" async></script>
    <link href="https://fonts.googleapis.com/css?family=Merriweather&display=swap" rel="stylesheet"> 
    <?php if (isset($jsHandler)) :?>
      <!-- Shared JavaScript is declared here. -->
      <script type="text/javascript">
       <?php if ($prod === true) :?>
       const prod = true;
       const apiGateway = '<?= $apiGateway . '/prod' ?>';
       const stripePublicKey = '<?= $stripeProdPublicKey ?>';
       <?php else :?>
       const prod = false;
       const apiGateway = '<?= $apiGateway . '/dev' ?>';
       const stripePublicKey = '<?= $stripeTestPublicKey ?>';
       <?php endif; ?>

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

       function setDisabled(selector, disabled) {
         const element = document.querySelector(selector);
         if (element && disabled) {
           element.setAttribute('disabled', true);
         } else if (element && !disabled) {
           element.removeAttribute('disabled');
         }
       }

       function formError(message, target) {
         const formErrors = document.getElementById('form-errors');
         if (!formErrors){
           console.error('expected #form-errors element');
           return
         }
         formErrors.textContent = message;

         const badInput = document.getElementById(target);
         if (!badInput) {
           return;
         }

         badInput.classList.add('form-error');
         badInput.addEventListener('change', () => {
           badInput.classList.remove('form-error');
         }, { once: true });
       }
      </script>
      <script type="text/javascript" src="<?=$jsHandler?>" async></script>
    <?php endif; ?>
    <title>Cosmo's Tuna</title>
  </head>
  <body>
    <header>
      <nav>
        <?php buildLinks() ?>
      </nav>
    </header>
