# Using the hosted fields feature with PayPlus.
PayPlus offers the ability to host your own payment page on your own website without relying on external redirect pages or iframes.
This repository contains a simple demo showing how to achieve that. 

In order to be able to run this demo, you need:
* Test/dev credentials from PayPlus, with API permission.
* A server with a domain name, an SSL certificate and support for PHP 7.4 or greater. The PHP part is technically optional because you could write your own code to accomplish this step in any server side language, but this demo only comes with a php file out of the box.
* Knowledge of npm package manager, html/javascript and basic php (again, the PHP part can be replaced with any other server-side language provided you can write your own code to create a new payment page link. You could even make this request in Postman and then inject the response into the browser).

The repository generally contains 3 files that are of importance for the demo: payment.php, index.html and script.js
* **payment.php** - this is where the original payment page request is generated. As payment page generation requires secret credentials, it is done on the website's server before this information can be passed on to the browser. This file merely contains the most rudimentary payment page generation request with some sample data. Do note, however, that you need to edit it to enter your own credentials and other information.
  ```php
    define('API_KEY', '');
    define('SECRET_KEY', '');
    define('PAYMENT_PAGE_UID', '');
    define('ORIGIN_DOMAIN', 'https://www.example.com');
    define('SUCCESS_URL', 'https://www.example.com/success');
    define('FAILURE_URL', 'https://www.example.com/failure');
    define('CANCEL_URL', 'https://www.example.com/cancel');
  ```
  As this demo mostly focuses on the hosted fields part, assuming you already know how to generate a payment page, I will leave it at that.
  However, there are some more useful comments inside the file itself.
* **index.html** - defines the payment page layout. This demo uses bootstrap for the design, just to make things more aesthetic.
* **script.js** - this is where the magic happens. Here we initialize the hosted fields plugin, map all the required fields, and define callbacks. Most of the demo happens here.
## Installation
* Download/clone this repository into a working server into whatever directory that is accessible via HTTP.
* Modify the file payment.php as explained above.
* Run ```npm i``` to install dependencies.
At this point, you should be able to navigate to the index of whatever directory you put this demo into.
You should see a payment form with some fields, and basic boostrap design.
  
## Overview
In a nutshell, the process starts when the client first makes a request to payment.php where a payment page link will be generated, and returned it to the client which then will parse the html file, replace placehodlers with iframes with actual hosted fields, take control over the rest of the non hosted fields, and will listen to the client interaction. 
Upon submission, it will send the information to the hosted field iframe which will in turn communicate with the PayPlus server to process the transaction, or return error. The client will then receive a response from the server via the hosted fields iframe, and will determine how to proceed further.

## script.js
We start by initializing the hosted fields dom object:
```javascript
const hf = new PayPlusHostedFieldsDom();
```
Then, we need to tell the plugin where the relevant fields are on the html page.
This is the part where we tell the plugin about the actual hosted fields. 
All of these are html tags that will either be replaced by an iframe containing the relevant field, or hidden if not needed.
We define the element selector, and also a wrapper element. This is to allow the plugin to hide the entire block if it needs to be hidden, and not just the field itself.

Also notice that we use simple CSS selectors. Any valid selector should work, both here and in the next section
```javascript
hf.SetMainFields({
	cc: {
		elmSelector: "#cc",
		wrapperElmSelector: "#cc-wrapper",
	},
	expiryy: {
		elmSelector: "#expiryy",
		wrapperElmSelector: ".expiry-wrapper",
	},
	expirym: {
		elmSelector: "#expirym",
		wrapperElmSelector: ".expiry-wrapper",
	},
	expiry: {
		elmSelector: "#expiry",
		wrapperElmSelector: ".expiry-wrapper-full",
	},
	cvv: {
		elmSelector: "#cvv",
		wrapperElmSelector: "#cvv-wrapper",
	},
})
```
Notice that we're defining 5 fields here. The first one is credit card. Self explanatory.
Second through fourth are for expiration date. The reason we need to define 3 of them is because PayPlus payment page configuration allows you to split up the expiry field into two fields (one for the month, one for the year).
So, the first two are for the two individual fields, and the last one - in case the page is configured to have a single expiry field. 
This way, if the page is configured to have separate expiry fields, this affords the website developer a more granular control over how the page will look regardless of the page configuration.

Also worth mentioning the CVV field. This is a pivotal field for the whole process. If payment page configuration requires showing it, it'll just show up as a field, otherwise, it will stay hidden but will nonetheless be still there, because this is the one field to rule them all.
This is the field that controls the whole process of communication between the hosted fields, and the payplus server, and the hosted fields and the client website. 
This is just for general knowledge as these inner workings happen automatically. Still, this is worth knowing.

```javascript
	.AddField("card_holder_id", "#id-number", "#id-number-wrapper")
	.AddField("payments", "#payments", "#payments-wrapper")
	.AddField("card_holder_name", "#card-holder-name", "#card-holder-name")
	.AddField("customer_name", "[name=customer_name]", ".customer_name-wrapper")
	.AddField("vat_number", "[name=customer_id]", ".customer_id-wrapper")
	.AddField("phone", "[name=phone]", ".phone-wrapper")
	.AddField("email", "[name=email]", ".email-wrapper")
	.AddField("contact_address", "[name=address]", ".address-wrapper")
	.AddField("contact_country", "[name=country]", ".country-wrapper")
	.AddField("custom_invoice_name", "#invoice-name", "#invoice-name-wrapper")
	.AddField("notes", "[name=notes]", ".notes-wrapper")
```
Here, similarly to the previous snippet, we define some other fields. These, unlike the first batch, will not be replaced by iframes but rather be hosted on under your domain's origin. 
We still need to map them for the plugin to be able to take control of them.
This is because upon form submission, this information will be sent over to the server as well. The plugin needs to know where these fields are so it can access them.
Also, if you choose to disable/hide them via the payment page configuration, the plugin will hide them as well. 
It is generally a wise idea to map as many of these fields as possible. Anything that isn't needed, will not be shown. On the other hand, should you choose to reconfigure the payment page to include them, they will already be there.
For instance, you may decide at first that you don't need the card_holder_name for now. But if you don't map it, and then decide to enable it via configuration, there will no field for it to use, so you would need to actually edit the page to add it manually. 
Best to map them all, just in case.

```javascript
$.get("payment.php", async (resp) => {
	if (resp.results.status == "success") {
		try {
			await hf.CreatePaymentPage({
				hosted_fields_uuid: resp.data.hosted_fields_uuid,
				page_request_uid: resp.data.page_request_uid,
				origin: 'https://restapidev.payplus.co.il'
			});
		} catch (error) {
			alert(error);
		}
		hf.InitPaymentPage.then((data) => {
			$("#create-payment-form").hide();
			$("#payment-form").show();
		});
	} else {
		alert(resp.results.message);
	}
});
```
This part will mostly differ from website to website depending on the framework used. I used simple jQuery for this entire demo as it is the easiest to understand, and most commonly used.
In a nutshell, what this code does is it first makes a GET request to the previously mentioned payment.php. Assuming a successful response, it will return 2 parameters that are needed for the process to start. We need the **page_request_uid**, and **hosted_fields_uuid** for the process to initialize.

As can be seen, CreatePaymentPage accepts an object with a third parameter, "origin". Not to be confused with the refURL_origin parameter in payment.php, this one is the url of the API server.
In this case I use **https://restapidev.payplus.co.il** for dev environment. For prod, it's **https://restapi.payplus.co.il**.

We then subscribe to the hosted field's InitPaymentPage promise, and upon its completion, we can do whatever housekeeping we may need to do (like actually showing the HTML form, hiding irrelevant stuff, whatever is needed by the website)

```javascript
$("#submit-payment").on("click", () => {
  hf.SubmitPayment();
});
```
Here I define the submit button. When the button is clicked, hf's SubmitPayment routine will commence. 

Lastly, we define some events
```javascript
hf.Upon("pp_pageExpired", (e) => {
	$("#submit-payment").prop("disabled", true);
	$("#status").val("Page Expired");
});

hf.Upon("pp_noAttemptedRemaining", (e) => {
	alert("No more attempts remaining");
});

hf.Upon("pp_responseFromServer", (e) => {
	let r = "";
	try {
		r = JSON.stringify(e.detail, null, 2);
	} catch (error) {
		r = e.detail;
	}
	$("#status").val(r);
});
hf.Upon("pp_submitProcess", (e) => {
	$("#submit-payment").prop("disabled", e.detail);
});
```
Using hf's method "Upon", we can subscribe to events and define custom behavior for when these events fire.
For the sake of this demo, I chose to use the pp_responseFromServer event to display whatever information server spits back onto a textarea for debugging purposes. 
In a real environment, this will probably be used for some sort of event routing like displaying a ThankYou page/popup upon success, or failure page upon failure.

## Fields
As mentioned above, the plugin expects 2 types of fields to be mapped.
Hosted fields require any element, as it will be replaced by an iframe anyway. Here I use a simple span tag.
Non hosted fields - actual fields that are entirely hosted on the website. These need to be either input or select fields. 
### Hosted fields:
* **cc** - Credit card field.
* **cvv** - Will naturally hold the field for CVV but is also the main field that controls the entire process.
* **expiry** - a single expiry field expecting the month and year of the credit card.
* **expiryy**/**expirym** - Expiry year/month fields respectively. Used when configuired to host separate fields for expiry year and month.
### Non hosted fields
* **card_holder_id** - ID/ Israeli national ID card number of the card owner
* **card_holder_name** - Name of the card owner
* **customer_name** - customer name, as defined internally on payplus or sent on the initial payment page request.
* **vat_number** - the Israeli ID card number, but similarly to the previous field, as it was submitted during payment page creation.
* **custom_invoice_name** - Alternative name to be used to issue a resulting invoice.
* **phone** - phone number.
* **email** - email address.
* **contact_address** - contact address.
* **contact_country** - contact country.
* **notes** - textual notes to be sent along with the transaction.
  
**Note:** you may wonder why there are two fields for the customer's name and their Israeli ID number. The first two fields strictly refer to the card holder information, whereas the next two refer to the customer. The two don't strictly have to be the same (think friend paying for another friend's bill ...etc.)

**Note:** there's more. PayPlus allows you to define any additional custom fields to be displayed on the payment page. The hosted fields plugin supports this as well.
  You can use the AddField method to add any other field, provided it was defined in the payment page settings.
  So, for instance, if your payment page has a custom field called "Code", you may add your own html field for it, and then simply map it, as you normally would, any other predefined hosted field:
  ```javascript
  .AddField('Code', '#code-fld', '#code-fld-wrapper')
  ```
  The hosted fields plugin will hide any mapped fields that aren't defined so feel free to define any possible fields, and they will only appear if/when configured.
  This will also handle their required/optional setting.

## Events
As mentioned above, the plugin exposes a number of events that will fire during various situations throughout the process.
To subscribe to an event, the plugin provides the function **Upon**:
```javascript
hf.Upon("pp_noAttemptedRemaining", (e) => {
	alert("No more attempts remaining");
});
```
2 parameters. Event name, and a callback with the event's data as a callback's sole parameter.
#### Available events:
* **pp_responseFromServer** - perhaps the most important one. Whenever there's a response from the server, this event will fire and will include its content in the callback's parameter. This is where you can ultimately check for errors returned by the server, or whether the request was successful.
* **pp_noAttemptedRemaining** - fires when the customer exceeds the number of available attempts.
* **pp_pageExpired** - when the payment page expires.
* **pp_paymentPageKilled** - when the payment page is invalidated. It fires along with the two events above
* **pp_submitProcess** - fires when information is submitted to the server for processing, and again when processing is over, indicating the status with either value true or false. Thus, when fired with "true", and until it fires again with "false", it can be assumed that the submit process is ongoing. Useful when you want to display a loader for the user indicating that the request is still being processed. For instance, in this example I use it to disable the submit button while processing, and reenable it when done.
* **pp_ccTypeChange** - when the customer enters their credit card number, and enough digits have been entered, the system is able to guess the brand of the credit card. For instance, if the digits 53 have been entered, the system assumes that it will be a Mastercard. This event fires whenever such detection occurs, usually when the customer entered the first two digits, but may happen multiple times if the customer clears the credit card field and types again. This is useful if you wish to display the credit card brand's logo as soon as it can be inferred. 
