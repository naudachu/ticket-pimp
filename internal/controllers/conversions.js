function myFunction() {
	var sheet = SpreadsheetApp.getActiveSheet();
	var data = sheet.getDataRange().getValues();
	var a1 = data[1][1];
	var a2 = ''
	var a3 = data[1][2];
	var reqUrl = data[1][3];
	var startarray = data[1][4]
	
	var regCount = 0;
	var depCount = 0;
	// reqUrl = 'https://graph.facebook.com/v15.0/841990173617973/activities?event=CUSTOM_APP_EVENTS&application_tracking_enabled=1&advertiser_tracking_enabled=1&advertiser_id=0a42bbee-c049-4180-ae3f-ce76ff33556e&841990173617973%7Cnjqit85G_uma0voikeP6eJvBiQk&custom_events=[%7B%22_eventName%22:%22fb_mobile_purchase%22%7D]'
	reqUrl = reqUrl.replace('A1', a1)
	reqUrl = reqUrl.replace('A3', a3)
	Logger.log("urll " + reqUrl)
	Logger.log("urll " + encodeURI(reqUrl).toString())
	
	for (var i = startarray; i < data.length; i++) {
	  a2 = data[i][0]
	  reqUrl = reqUrl.replace('A2', a2)
  
	  var params = {
		"method": 'post'
	  }
	  var response = UrlFetchApp.fetch(encodeURI(reqUrl), params);
	  Logger.log(i + " " + response.getContentText() + ' ' + response.getResponseCode());
	  depCount++;
  
	}
	Logger.log('Reg Count= ' + regCount);
	Logger.log('Dep Count= ' + depCount);
  
  }
