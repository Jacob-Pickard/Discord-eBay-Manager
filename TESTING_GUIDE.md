# Quick Testing Guide - Sandbox Offer Management

## 🎯 Goal
Test all offer management features before production deployment.

## ⚡ Quick Start (15 minutes)

### Step 1: Create Test Listing (5 min)
1. Go to: https://sandbox.ebay.com
2. Sign in with your sandbox seller account
3. Click "Sell" → "Create listing"
4. Fill in:
   - Title: "Test Item - Vintage Watch"
   - Category: Jewelry & Watches
   - Price: $99.99
   - ✅ **Enable "Best Offer"**
   - Min acceptable offer: $75.00
5. Click "List item"
6. **Copy the Item ID** from the URL

### Step 2: Create Sandbox Buyer (3 min)
1. Open incognito/private browser window
2. Go to: https://developer.ebay.com/sandbox/register
3. Create a new test buyer account
4. Verify email (check sandbox inbox)

### Step 3: Make Test Offer (2 min)
1. In incognito window, search for your listing
2. Click "Make Offer"
3. Enter offer: $60.00
4. Submit offer
5. **Note the offer will appear in seller account**

### Step 4: Test Discord Commands (5 min)

#### 4a. View Offers
```
/get-offers
```
**Expected:** Shows your pending offer with details

#### 4b. Test Counter Offer
```
/counter-offer offer-id:<paste_id> price:70.00
```
**Expected:** 
- ✅ Success message
- Buyer receives counter in their account

#### 4c. Check Buyer Response (in incognito)
1. Log in as buyer
2. Go to "My eBay" → "Offers"
3. See your counter offer
4. Can accept or counter again

#### 4d. Test Accept
Make another offer at $80, then:
```
/accept-offer offer-id:<new_offer_id>
```
**Expected:**
- ✅ Creates order
- Shows in `/get-orders`
- Buyer gets purchase confirmation

#### 4e. Test Decline
Make another offer at $50, then:
```
/decline-offer offer-id:<low_offer_id>
```
**Expected:**
- ✅ Offer closed
- Buyer notified
- Removed from pending

---

## 🧪 Detailed Test Cases

### Test Case 1: Counter Offer Flow
**Steps:**
1. Buyer offers $60 on $99.99 item
2. Seller counters at $80
3. Buyer accepts counter

**Commands:**
```bash
# View the offer
/get-offers

# Counter it
/counter-offer offer-id:ABC123 price:80.00

# (Buyer accepts in browser)

# Verify order created
/get-orders
```

**Expected Results:**
- Counter sent successfully
- Buyer receives notification
- Order created when buyer accepts
- Webhook notification received

---

### Test Case 2: Accept Offer
**Steps:**
1. Buyer offers $85 on $99.99 item (good offer)
2. Seller accepts immediately

**Commands:**
```bash
/get-offers
/accept-offer offer-id:XYZ789
/get-orders
```

**Expected Results:**
- Offer accepted successfully
- Order created immediately
- Order shows PAID status
- Both parties receive confirmation

---

### Test Case 3: Decline Offer
**Steps:**
1. Buyer offers $30 on $99.99 item (too low)
2. Seller declines

**Commands:**
```bash
/get-offers
/decline-offer offer-id:LOW123
/get-offers # Should not show declined offer
```

**Expected Results:**
- Offer declined successfully
- Buyer receives decline notification
- Offer removed from pending list

---

### Test Case 4: Multiple Offers
**Steps:**
1. Create 3 offers from different amounts
2. Counter one
3. Accept one
4. Decline one

**Commands:**
```bash
# View all offers
/get-offers

# Counter first (Offer $60)
/counter-offer offer-id:OFFER1 price:75.00

# Accept second (Offer $85)
/accept-offer offer-id:OFFER2

# Decline third (Offer $40)
/decline-offer offer-id:OFFER3

# Verify results
/get-offers  # Should only show counter-offer pending
/get-orders  # Should show accepted offer as order
```

---

### Test Case 5: Error Handling
**Test invalid inputs to verify error messages:**

```bash
# Invalid offer ID
/accept-offer offer-id:FAKE123
# Expected: "Failed to accept offer" with troubleshooting

# Invalid price format
/counter-offer offer-id:REAL123 price:abc
# Expected: "Invalid price format"

# Negative price
/counter-offer offer-id:REAL123 price:-50
# Expected: "Price must be greater than $0.00"

# Expired offer (process same offer twice)
/accept-offer offer-id:SAME123
/accept-offer offer-id:SAME123
# Expected: Second attempt fails with error
```

---

## 📊 Test Results Checklist

### ✅ Offer Viewing
- [ ] `/get-offers` shows pending offers
- [ ] Displays offer amount correctly
- [ ] Shows buyer information
- [ ] Shows item details
- [ ] Shows offer timestamp

### ✅ Counter Offer
- [ ] `/counter-offer` accepts valid offer ID
- [ ] Validates price format
- [ ] Rejects negative prices
- [ ] Sends counter to eBay API
- [ ] Shows success message
- [ ] Buyer receives counter notification
- [ ] Buyer can see counter in their account

### ✅ Accept Offer
- [ ] `/accept-offer` accepts valid offer ID
- [ ] Creates order on eBay
- [ ] Order appears in `/get-orders`
- [ ] Success message displays
- [ ] Buyer receives confirmation
- [ ] Webhook notification received (if configured)

### ✅ Decline Offer
- [ ] `/decline-offer` accepts valid offer ID
- [ ] Removes offer from pending
- [ ] Shows success message
- [ ] Buyer receives decline notification
- [ ] Offer no longer appears in `/get-offers`

### ✅ Error Handling
- [ ] Invalid offer ID shows error
- [ ] Invalid price format shows error
- [ ] Expired offer shows appropriate error
- [ ] Authorization errors handled gracefully
- [ ] Error messages include troubleshooting tips

---

## 🐛 Common Issues & Solutions

### Issue: "No offers found"
**Cause:** Listing doesn't have Best Offer enabled  
**Solution:** Edit listing, enable "Best Offer" in pricing section

### Issue: Can't find offer ID
**Cause:** Offer ID not displayed clearly  
**Solution:** Check `/get-offers` - ID is in the response. Look for pattern like "ABC123..."

### Issue: "Failed to counter offer"
**Cause:** Offer already expired or accepted  
**Solution:** Check offer status in eBay Seller Hub

### Issue: Buyer doesn't see counter
**Cause:** Buyer notification delay  
**Solution:** Wait 1-2 minutes, refresh buyer's "Offers" page

### Issue: Order not showing after accept
**Cause:** Order processing delay  
**Solution:** Wait 30 seconds, run `/get-orders` again

---

## 📝 Test Report Template

After completing tests, fill this out:

```
## Test Report - Sandbox Offer Management
Date: ___________
Tester: ___________

### Environment
- Bot Version: v1.0
- eBay Environment: Sandbox
- Discord Server: ___________

### Test Results
✅ `/get-offers` - PASS / FAIL
   Notes: _______________________

✅ `/counter-offer` - PASS / FAIL
   Notes: _______________________

✅ `/accept-offer` - PASS / FAIL
   Notes: _______________________

✅ `/decline-offer` - PASS / FAIL
   Notes: _______________________

✅ Error Handling - PASS / FAIL
   Notes: _______________________

### Issues Found
1. _______________________
2. _______________________
3. _______________________

### Recommendations
- Ready for production: YES / NO
- Additional testing needed: _______________________
- Concerns: _______________________

### Sign-Off
Tested by: ___________
Date: ___________
Approved for Production: YES / NO
```

---

## ✅ Success Criteria

Before moving to production, verify:

- [x] All offer commands work without errors
- [x] Orders created successfully from accepted offers
- [x] Error messages are clear and helpful
- [x] Buyer receives all notifications
- [x] Webhook notifications arrive (if configured)
- [x] No bot crashes during testing
- [x] Response times are reasonable (< 5 seconds)

**Once all checked, you're ready for production! 🚀**

---

## 🎓 Pro Tips

1. **Use Descriptive Titles:** Make test listings easy to identify
2. **Test Multiple Scenarios:** Don't just test happy path
3. **Check Both Sides:** Verify buyer and seller perspectives
4. **Document Everything:** Screenshot errors for troubleshooting
5. **Test Edge Cases:** Expired offers, duplicate actions, etc.
6. **Time Your Tests:** Note how long each operation takes

---

## 📞 Need Help?

If you encounter issues during testing:

1. Check bot-error.log for detailed errors
2. Run `/ebay-status` to verify authorization
3. Check PRODUCTION_READINESS.md troubleshooting section
4. Review eBay API documentation for specific endpoints
5. Verify sandbox account is properly configured

**Testing complete? See PRODUCTION_READINESS.md for deployment steps!**
