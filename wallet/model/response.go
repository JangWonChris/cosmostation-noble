package model

/*
	API 호출 후 리턴되는 리턴값을 동일한 포맷으로 리턴해주고 싶다.
	성공적으로 호출이 됬을 경우,
	{
		result: "success",
		data: [
			alarm_token: "ADFKJAEMNQDASF",
			device_type: "android",
			coin_type: "ATOM",
			status: true,
		]
	}

	실패했을 경우,
	{
		error_code: 202,
		error_msg: "Duplicate account and address",
	}
*/

type AccountResponse struct {
	Result string    `json:"result"`
	Data   []Account `json:"data"`
}
