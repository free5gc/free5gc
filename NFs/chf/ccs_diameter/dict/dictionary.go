package dict

import "encoding/xml"

// Only list the AVP details are needed for the class B method for ServiceUsage in TS 32.296
const (
	RateDictionary = xml.Header + `
	<diameter>
	<application id="16777218">
		<command code="111" short="SU" name="Service-Usage">
		<request>
			<rule avp="Session-Id" required="true" max="1"/>
			<rule avp="Origin-Host" required="true" max="1"/>
			<rule avp="Origin-Realm" required="true" max="1"/>
			<rule avp="Destination-Realm" required="true" max="1"/>
			<rule avp="Destination-Host" required="true" max="1"/>
			<rule avp="Vendor-Specific-Application-Id" required="false" max="1"/>
			<rule avp="User-Name" required="false" max="1"/>
			<rule avp="Event-Timestamp" required="false" max="1"/>
			<rule avp="BeginTime" required="false" max="1"/>
			<rule avp="ActualTime" required="false" max="1"/>
			<rule avp="Subscription-Id" required="false" max="1"/>
			<rule avp="Service-Rating" required="false" max="1"/>
		</request>
		<answer>
			<rule avp="Session-Id" required="true" max="1"/>
			<rule avp="Origin-Host" required="true" max="1"/>
			<rule avp="Origin-Realm" required="true" max="1"/>
			<rule avp="Vendor-Specific-Application-Id" required="false" max="1"/>
			<rule avp="Event-Timestamp" required="false" max="1"/>
			<rule avp="Service-Rating" required="false" max="1"/>
		</answer>
		</command>

		<avp name="BeginTime" code="7000">
			<data type="Time"/>
		</avp>

		<avp name="ActualTime" code="7001">
			<data type="Time"/>
		</avp>

		<avp name="Service-Rating" code="7002">
			<data type="Grouped">
				<rule avp="Service-Identifier" required="true" max="1"/>
				<rule avp="DestinationID" required="false" max="1"/>
				<rule avp="ServiceInformation" required="false" max="1"/>
				<rule avp="Extension" required="false" max="1"/>
				<rule avp="Price" required="false" max="1"/>
				<rule avp="BillingInfo" required="false" max="1"/>
				<rule avp="TariffSwitchTime" required="false" max="1"/>
				<rule avp="MonetaryTariff" required="true" max="1"/>
				<rule avp="NextMonetaryTariff" required="false" max="1"/>
				<rule avp="ExpiryTime" required="false" max="1"/>
				<rule avp="ValidUnits" required="false" max="1"/>
				<rule avp="MonetaryTariffAfterValidUnits" required="false" max="1"/>
				<rule avp="Counter" required="false" max="1"/>
				<rule avp="BasicPriceTimeStamp" required="false" max="1"/>
				<rule avp="BasicPrice" required="false" max="1"/>
				<rule avp="CounterPrice" required="false" max="1"/>
				<rule avp="CounterTariff" required="false" max="1"/>
				<rule avp="RequestedCounters" required="false" max="1"/>
				<rule avp="RequestSubType" required="false" max="1"/>
				<rule avp="ImpactonCounter" required="false" max="1"/>
				<rule avp="RequestedUnits" required="false" max="1"/>
				<rule avp="ConsumedUnits" required="false" max="1"/>
				<rule avp="ConsumedUnitsAfterTariffSwitch" required="false" max="1"/>
				<rule avp="MonetaryQuota" required="false" max="1"/>
				<rule avp="MinimalRequestedUnits" required="false" max="1"/>
				<rule avp="AllowedUnits" required="false" max="1"/>
			</data>
	</avp>

		<avp name="DestinationID" code="7003">
			<data type="Grouped">
				<rule avp="DestinationIDType" required="true" max="1"/>
				<rule avp="DestinationIDData" required="true" max="1"/>
			</data>
		</avp>

		<avp name="Extension" code="7004">
			<data type="Grouped"/>
		</avp>

		<avp name="Price" code="7005">
			<data type="Unsigned32"/>
		</avp>

		<avp name="BillingInfo" code="7006">
			<data type="UTF8String"/>
		</avp>

		<avp name="TariffSwitchTime" code="7007">
			<data type="Unsigned32"/>
		</avp>

		<avp name="MonetaryTariff" code="7008">
			<data type="Grouped">
				<rule avp="Currency-Code" required="false" max="1"/>
				<rule avp="Scale-Factor" required="false" max="1"/>
				<rule avp="Rate-Element" required="false" max="1"/>
			</data>
		</avp>

		<avp name="NextMonetaryTariff" code="7009">
			<data type="Grouped">
				<rule avp="Currency-Code" required="false" max="1"/>
				<rule avp="Scale-Factor" required="false" max="1"/>
				<rule avp="Rate-Element" required="false" max="1"/>
			</data>
		</avp>

		<avp name="ExpiryTime" code="7010">
			<data type="Time"/>
		</avp>

		<avp name="ValidUnits" code="7011">
			<data type="Unsigned32"/>
		</avp>

		<avp name="MonetaryTariffAfterValidUnits" code="7012">
			<data type="Grouped">
				<rule avp="Currency-Code" required="false" max="1"/>
				<rule avp="Scale-Factor" required="false" max="1"/>
				<rule avp="Rate-Element" required="false" max="1"/>
			</data>
		</avp>

		<avp name="RequestSubType" code="7013">
			<data type="Enumerated">
				<item code="0" name="REQ_SUBTYPE_AOC"/>
				<item code="1" name="REQ_SUBTYPE_RESERVE"/>
				<item code="2" name="REQ_SUBTYPE_DEBIT"/>
				<item code="3" name="REQ_SUBTYPE_RELEASE"/>
			</data>
		</avp>

		<avp name="ConsumedUnits" code="7014">
			<data type="Unsigned32"/>
		</avp>

		<avp name="ConsumedUnitsAfterTariffSwitch" code="7015">
			<data type="Unsigned32"/>
		</avp>

		<avp name="MonetaryQuota" code="7016">
			<data type="Unsigned32"/>
		</avp>

		<avp name="RequestedUnits" code="7017">
			<data type="Unsigned32"/>
		</avp>

		<avp name="MinimalRequestedUnits" code="7018">
			<data type="Unsigned32"/>
		</avp>

		<avp name="ServiceInformation" code="7019">
			<data type="Grouped"/>
		</avp>

		<avp name="ImpactonCounter" code="7020">
			<data type="Grouped"/>
		</avp>

		<avp name="AllowedUnits" code="7021">
			<data type="Unsigned32"/>
		</avp>

		<avp name="RequestedCounters" code="7022">
			<data type="Grouped"/>
		</avp>

		<avp name="CounterTariff" code="7023">
			<data type="Grouped"/>
		</avp>

		<avp name="CounterPrice" code="7024">
			<data type="Grouped"/>
		</avp>

		<avp name="BasicPriceTimeStamp" code="7025">
			<data type="Time"/>
		</avp>

		<avp name="Counter" code="7026">
			<data type="Grouped"/>
		</avp>

		<avp name="Vendor-Specific-Application-Id" code="7027">
			<data type="Grouped"/>
		</avp>

		<avp name="Currency-Code" code="425" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.11 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Subscription-Id" code="443" must="M" may="P" must-not="V" may-encrypt="Y" vendor-id="0">
			<!-- https://tools.ietf.org/rfc/rfc4006.txt -->
			<data type="Grouped">
				<rule avp="Subscription-Id-Type" required="false" max="1"/>
				<rule avp="Subscription-Id-Data" required="false" max="1"/>
			</data>
		</avp>

		<avp name="Subscription-Id-Type" code="450" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.47-->
			<data type="Enumerated">
				<item code="0" name="END_USER_E164"/>
				<item code="1" name="END_USER_IMSI"/>
				<item code="2" name="END_USER_SIP_URI"/>
				<item code="3" name="END_USER_NAI"/>
			</data>
		</avp>

		<avp name="Service-Identifier" code="439" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.28-->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Subscription-Id-Data" code="444" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.48-->
			<data type="UTF8String"/>
		</avp>

		<avp name="Rate-Element" code="2058" must="V,M" may="P" must-not="-" may-encrypt="N" >
			<data type="Grouped">
				<rule avp="CC-Unit-Type" required="true" max="1"/>
				<rule avp="Charge-Reason-Code" required="false" max="1"/>
				<rule avp="Unit-Value" required="false" max="1"/>
				<rule avp="Unit-Cost" required="false" max="1"/>
				<rule avp="Unit-Quota-Threshold" required="false" max="1"/>
			</data>
		</avp>

		<avp name="Unit-Quota-Threshold" code="1226" must="V,M" may="P" must-not="-" may-encrypt="N" vendor-id="10415">
			<data type="Unsigned32"/>
		</avp>
		
		<avp name="Scale-Factor" code="2059" must="V,M" may="P" must-not="-" may-encrypt="N">
			<data type="Grouped">
				<rule avp="Value-Digits" required="true" max="1"/>
				<rule avp="Exponent" required="false" max="1"/>
			</data>
		</avp>

		<avp name="Value-Digits" code="447" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.10-->
			<data type="Integer64"/>
		</avp>

		<avp name="Unit-Value" code="445" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.8-->
			<data type="Grouped">
				<rule avp="Value-Digits" required="true" max="1"/>
				<rule avp="Exponent" required="true" max="1"/>
			</data>
		</avp>

		<avp name="Exponent" code="429" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.9 -->
			<data type="Integer32"/>
		</avp>

		<avp name="CC-Unit-Type" code="454" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.32 -->
			<data type="Enumerated">
				<item code="0" name="TIME"/>
				<item code="1" name="MONEY"/>
				<item code="2" name="TOTAL-OCTETS"/>
				<item code="3" name="INPUT-OCTETS"/>
				<item code="4" name="OUTPUT-OCTETS"/>
				<item code="5" name="SERVICE-SPECIFIC-UNITS"/>
			</data>
		</avp>

		<avp name="Charge-Reason-Code" code="2118" must="V,M" may="P" must-not="-" may-encrypt="N">
			<data type="Enumerated">
				<item code="0" name="UNKNOWN"/>
				<item code="1" name="USAGE"/>
				<item code="2" name="COMMUNICATION-ATTEMPT-CHARGE"/>
				<item code="3" name="SETUP-CHARGE"/>
				<item code="4" name="ADD-ON-CHARGE"/>
			</data>
		</avp>

		<avp name="Unit-Cost" code="2061" must="V,M" may="P" must-not="-" may-encrypt="N">
			<data type="Grouped">
				<rule avp="Value-Digits" required="true" max="1"/>
				<rule avp="Exponent" required="false" max="1"/>
			</data>
		</avp>

	</application>
	</diameter>
	`

	// Copy from RFC 4006, but some avp are added or deleted according to 32.296 B.5
	AbmfDictionary = xml.Header + `
	<diameter>

	<application id="16777218" type="auth" name="Charging Control">
		<!-- Diameter Credit Control Application -->
		<!-- http://tools.ietf.org/html/rfc4006 -->

		<command code="272" short="CC" name="Credit-Control">
			<request>
				<!-- http://tools.ietf.org/html/rfc4006#section-3.1 -->
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="Destination-Realm" required="true" max="1"/>
				<rule avp="Auth-Application-Id" required="true" max="1"/>
				<rule avp="Service-Context-Id" required="true" max="1"/>
				<rule avp="CC-Request-Type" required="true" max="1"/>
				<rule avp="CC-Request-Number" required="true" max="1"/>
				<rule avp="Origin-State-Id" required="false" max="1"/>
				<rule avp="Destination-Host" required="false" max="1"/>
				<rule avp="User-Name" required="false" max="1"/>
				<rule avp="Event-Timestamp" required="false" max="1"/>
				<rule avp="Subscription-Id" required="false" max="1"/>
				<rule avp="Termination-Cause" required="false" max="1"/>
				<rule avp="Service-Identifier" required="false" max="1"/>
				<rule avp="Requested-Action" required="false" max="1"/>
				<rule avp="Multiple-Services-Indicator" required="false" max="1"/>
				<rule avp="Multiple-Services-Credit-Control" required="false" max="1"/>
				<rule avp="Proxy-Info" required="false" max="1"/>
				<rule avp="Service-Information" required="false" max="1"/>
			</request>
			<answer>
				<!-- http://tools.ietf.org/html/rfc4006#section-3.2 -->
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Result-Code" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="CC-Request-Type" required="true" max="1"/>
				<rule avp="CC-Request-Number" required="true" max="1"/>
				<rule avp="CC-Session-Failover" required="false" max="1"/>
				<rule avp="Multiple-Services-Credit-Control" required="false" max="1"/>
				<rule avp="Cost-Information" required="false" max="1"/>
				<rule avp="Low-Balance-Indication" required="false" max="1"/>
				<rule avp="Remaining-Balance" required="false" max="1"/>
				<rule avp="AB-Response" required="false" max="1"/>
				<rule avp="Credit-Control-Failure-Handling" required="false" max="1"/>
				<rule avp="Direct-Debiting-Failure-Handling" required="false" max="1"/>
				<rule avp="Validity-Time" required="false" max="1"/>
				<rule avp="Redirect-Host" required="false" max="1"/>
				<rule avp="Redirect-Host-Usage" required="false" max="1"/>
				<rule avp="Redirect-Max-Cache-Time" required="false" max="1"/>
				<rule avp="Proxy-Info" required="false" max="1"/>
				<rule avp="Failed-AVP" required="false" max="1"/>
				<rule avp="Service-Information" required="false" max="1"/>
			</answer>
		</command>

		<avp name="CC-Correlation-Id" code="411" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.1 -->
			<data type="OctetString"/>
		</avp>

		<avp name="CC-Input-Octets" code="412" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.24 -->
			<data type="Unsigned64"/>
		</avp>

		<avp name="CC-Money" code="413" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.22 -->
			<data type="Grouped">
				<rule avp="Unit-Value" required="true" max="1"/>
				<rule avp="Currency-Code" required="true" max="1"/>
			</data>
		</avp>

		<avp name="CC-Output-Octets" code="414" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.25 -->
			<data type="Unsigned64"/>
		</avp>

		<avp name="CC-Request-Number" code="415" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.2 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="CC-Request-Type" code="416" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.3 -->
			<data type="Enumerated">
				<item code="1" name="INITIAL_REQUEST"/>
				<item code="2" name="UPDATE_REQUEST"/>
				<item code="3" name="TERMINATION_REQUEST"/>
			</data>
		</avp>

		<avp name="CC-Service-Specific-Units" code="417" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Unsigned64"/>
		</avp>

		<avp name="CC-Session-Failover" code="418" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.4 -->
			<data type="Enumerated">
				<item code="0" name="FAILOVER_NOT_SUPPORTED"/>
				<item code="1" name="FAILOVER_SUPPORTED"/>
			</data>
		</avp>

		<avp name="CC-Sub-Session-Id" code="419" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.5 -->
			<data type="Unsigned64"/>
		</avp>

		<avp name="CC-Time" code="420" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.21 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="CC-Total-Octets" code="421" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.23 -->
			<data type="Unsigned64"/>
		</avp>

		<avp name="CC-Unit-Type" code="454" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.32 -->
			<data type="Enumerated">
				<item code="0" name="TIME"/>
				<item code="1" name="MONEY"/>
				<item code="2" name="TOTAL-OCTETS"/>
				<item code="3" name="INPUT-OCTETS"/>
				<item code="4" name="OUTPUT-OCTETS"/>
				<item code="5" name="SERVICE-SPECIFIC-UNITS"/>
			</data>
		</avp>

		<avp name="Check-Balance-Result" code="422" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.6 -->
			<data type="Enumerated">
				<item code="0" name="ENOUGH_CREDIT"/>
				<item code="1" name="NO_CREDIT"/>
			</data>
		</avp>

		<avp name="Cost-Information" code="423" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.7 -->
			<data type="Grouped">
				<rule avp="Unit-Value" required="true" max="1"/>
				<rule avp="Currency-Code" required="true" max="1"/>
				<rule avp="Cost-Unit" required="true" max="1"/>
			</data>
		</avp>

		<avp name="Cost-Unit" code="424" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.12 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Credit-Control" code="426" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.13 -->
			<data type="Enumerated">
				<item code="0" name="CREDIT_AUTHORIZATION"/>
				<item code="1" name="RE_AUTHORIZATION"/>
			</data>
		</avp>

		<avp name="Credit-Control-Failure-Handling" code="427" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.14 -->
			<data type="Enumerated">
				<item code="0" name="TERMINATE"/>
				<item code="1" name="CONTINUE"/>
				<item code="2" name="RETRY_AND_TERMINATE"/>
			</data>
		</avp>

		<avp name="Currency-Code" code="425" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.11 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Direct-Debiting-Failure-Handling" code="428" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.15 -->
			<data type="Enumerated">
				<item code="0" name="TERMINATE_OR_BUFFER"/>
				<item code="1" name="CONTINUE"/>
			</data>
		</avp>

		<avp name="Exponent" code="429" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.9 -->
			<data type="Integer32"/>
		</avp>

		<avp name="Final-Unit-Action" code="449" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.35 -->
			<data type="Enumerated">
				<item code="0" name="TERMINATE"/>
				<item code="1" name="REDIRECT"/>
				<item code="2" name="RESTRICT_ACCESS"/>
			</data>
		</avp>

		<avp name="Final-Unit-Indication" code="430" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.34 -->
			<data type="Grouped">
				<rule avp="Final-Unit-Action" required="true" max="1"/>
				<rule avp="Restriction-Filter-Rule" required="false" max="1"/>
				<rule avp="Filter-Id" required="false" max="1"/>
				<rule avp="Redirect-Server" required="false" max="1"/>
			</data>
		</avp>

		<avp name="Granted-Service-Unit" code="431" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.17 -->
			<data type="Grouped">
				<rule avp="Tariff-Time-Change" required="false" max="1"/>
				<rule avp="CC-Time" required="false" max="1"/>
				<rule avp="CC-Money" required="false" max="1"/>
				<rule avp="CC-Total-Octets" required="false" max="1"/>
				<rule avp="CC-Input-Octets" required="false" max="1"/>
				<rule avp="CC-Output-Octets" required="false" max="1"/>
				<rule avp="CC-Service-Specific-Units" required="false" max="1"/>
				<!-- *[ AVP ]-->
			</data>
		</avp>

		<avp name="G-S-U-Pool-Identifier" code="453" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.31 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="G-S-U-Pool-Reference" code="457" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.30 -->
			<data type="Grouped">
				<rule avp="G-S-U-Pool-Identifier" required="true" max="1"/>
				<rule avp="CC-Unit-Type" required="true" max="1"/>
				<rule avp="Unit-Value" required="true" max="1"/>
			</data>
		</avp>

		<avp name="Multiple-Services-Credit-Control" code="456" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.16 -->
			<data type="Grouped">
				<rule avp="Granted-Service-Unit" required="false" max="1"/>
				<rule avp="Requested-Service-Unit" required="false" max="1"/>
				<rule avp="Used-Service-Unit" required="false" max="1"/>
				<rule avp="Tariff-Change-Usage" required="false" max="1"/>
				<rule avp="Service-Identifier" required="false" max="1"/>
				<rule avp="Rating-Group" required="false" max="1"/>
				<rule avp="G-S-U-Pool-Reference" required="false" max="1"/>
				<rule avp="Validity-Time" required="false" max="1"/>
				<rule avp="Result-Code" required="false" max="1"/>
				<rule avp="Final-Unit-Indication" required="false" max="1"/>
				<!-- *[ AVP ]-->
			</data>
		</avp>

		<avp name="Multiple-Services-Indicator" code="455" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.40 -->
			<data type="Enumerated">
				<item code="0" name="MULTIPLE_SERVICES_NOT_SUPPORTED"/>
				<item code="1" name="MULTIPLE_SERVICES_SUPPORTED"/>
			</data>
		</avp>

		<avp name="Rating-Group" code="432" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.29 -->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Redirect-Address-Type" code="433" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.38 -->
			<data type="Enumerated">
				<item code="0" name="IPv4 Address"/>
				<item code="1" name="IPv6 Address"/>
				<item code="2" name="URL"/>
				<item code="3" name="SIP URI"/>
			</data>
		</avp>

		<avp name="Redirect-Server" code="434" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.37 -->
			<data type="Grouped">
				<rule avp="Redirect-Address-Type" required="true" max="1"/>
				<rule avp="Redirect-Server-Address" required="true" max="1"/>
			</data>
		</avp>

		<avp name="Redirect-Server-Address" code="435" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.39 -->
			<data type="UTF8String"/>
		</avp>

		<avp name="Requested-Action" code="436" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.41 -->
			<data type="Enumerated">
				<item code="0" name="DIRECT_DEBITING"/>
				<item code="1" name="REFUND_ACCOUNT"/>
				<item code="2" name="CHECK_BALANCE"/>
				<item code="3" name="PRICE_ENQUIRY"/>
			</data>
		</avp>

		<avp name="Requested-Service-Unit" code="437" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.18-->
			<data type="Grouped">
				<rule avp="CC-Time" required="false" max="1"/>
				<rule avp="CC-Money" required="false" max="1"/>
				<rule avp="CC-Total-Octets" required="false" max="1"/>
				<rule avp="CC-Input-Octets" required="false" max="1"/>
				<rule avp="CC-Output-Octets" required="false" max="1"/>
				<rule avp="CC-Service-Specific-Units" required="false" max="1"/>
				<!-- *[ AVP ]-->
			</data>
		</avp>

		<avp name="Restriction-Filter-Rule" code="438" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.36-->
			<data type="IPFilterRule"/>
		</avp>

		<avp name="Service-Context-Id" code="461" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.42-->
			<data type="UTF8String"/>
		</avp>

		<avp name="Service-Identifier" code="439" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.28-->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Service-Parameter-Info" code="440" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.43-->
			<data type="Grouped">
				<rule avp="Service-Parameter-Type" required="true" max="1"/>
				<rule avp="Service-Parameter-Value" required="true" max="1"/>
			</data>
		</avp>

		<avp name="Service-Parameter-Type" code="441" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.44-->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Service-Parameter-Value" code="442" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.45-->
			<data type="OctetString"/>
		</avp>

		<avp name="Subscription-Id" code="443" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.46-->
			<data type="Grouped">
				<rule avp="Subscription-Id-Type" required="true" max="1"/>
				<rule avp="Subscription-Id-Data" required="true" max="1"/>
			</data>
		</avp>

		<avp name="Subscription-Id-Data" code="444" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.48-->
			<data type="UTF8String"/>
		</avp>

		<avp name="Subscription-Id-Type" code="450" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.47-->
			<data type="Enumerated">
				<item code="0" name="END_USER_E164"/>
				<item code="1" name="END_USER_IMSI"/>
				<item code="2" name="END_USER_SIP_URI"/>
				<item code="3" name="END_USER_NAI"/>
			</data>
		</avp>

		<avp name="Tariff-Change-Usage" code="452" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.27-->
			<data type="Enumerated">
				<item code="0" name="UNIT_BEFORE_TARIFF_CHANGE"/>
				<item code="1" name="UNIT_AFTER_TARIFF_CHANGE"/>
				<item code="2" name="UNIT_INDETERMINATE"/>
			</data>
		</avp>

		<avp name="Tariff-Time-Change" code="451" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.20-->
			<data type="Time"/>
		</avp>

		<avp name="Unit-Value" code="445" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.8-->
			<data type="Grouped">
				<rule avp="Value-Digits" required="true" max="1"/>
				<rule avp="Exponent" required="true" max="1"/>
			</data>
		</avp>

		<avp name="Used-Service-Unit" code="446" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.19-->
			<data type="Grouped">
				<rule avp="Tariff-Change-Usage" required="false" max="1"/>
				<rule avp="CC-Time" required="false" max="1"/>
				<rule avp="CC-Money" required="false" max="1"/>
				<rule avp="CC-Total-Octets" required="false" max="1"/>
				<rule avp="CC-Input-Octets" required="false" max="1"/>
				<rule avp="CC-Output-Octets" required="false" max="1"/>
				<rule avp="CC-Service-Specific-Units" required="false" max="1"/>
				<!-- *[ AVP ]-->
			</data>
		</avp>

		<avp name="User-Equipment-Info" code="458" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.49-->
			<data type="Grouped">
				<rule avp="User-Equipment-Info-Type" required="true" max="1"/>
				<rule avp="User-Equipment-Info-Value" required="true" max="1"/>
			</data>
		</avp>

		<avp name="User-Equipment-Info-Type" code="459" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.50-->
			<data type="Enumerated">
				<item code="0" name="IMEISV"/>
				<item code="1" name="MAC"/>
				<item code="2" name="EUI64"/>
				<item code="3" name="MODIFIED_EUI64"/>
			</data>
		</avp>

		<avp name="User-Equipment-Info-Value" code="460" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.51-->
			<data type="OctetString"/>
		</avp>

		<avp name="Value-Digits" code="447" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.10-->
			<data type="Integer64"/>
		</avp>

		<avp name="Validity-Time" code="448" must="M" may="P" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc4006#section-8.33-->
			<data type="Unsigned32"/>
		</avp>

		<avp name="Service-Information" code="873" must="V,M" may="P" must-not="-" may-encrypt="N" vendor-id="10415">
			<data type="Grouped">
				<rule avp="Subscription-Id" required="false"/>
				<rule avp="AoC-Information" required="false" max="1"/>
				<rule avp="PS-Information" required="false" max="1"/>
				<rule avp="IMS-Information" required="false" max="1"/>
				<rule avp="MMS-Information" required="false" max="1"/>
				<rule avp="LCS-Information" required="false" max="1"/>
				<rule avp="PoC-Information" required="false" max="1"/>
				<rule avp="MBMS-Information" required="false" max="1"/>
				<rule avp="SMS-Information" required="false" max="1"/>
				<rule avp="VCS-Information" required="false" max="1"/>
				<rule avp="MMTel-Information" required="false" max="1"/>
				<rule avp="Service-Generic-Information" required="false" max="1"/>
				<rule avp="IM-Information" required="false" max="1"/>
				<rule avp="DCD-Information" required="false" max="1"/>
			</data>
		</avp>

		<avp name="Filter-Id" code="11" must="M" may="" must-not="V" may-encrypt="Y">
			<!-- http://tools.ietf.org/html/rfc7155#section-4.4.7 -->
			<data type="UTF8String"/>
		</avp>


		<avp name="Low-Balance-Indication" code="2020">
			<data type="Enumerated">
				<item code="0" name="NOT-APPLICABLE"/>
				<item code="1" name="YES"/>
			</data>
		</avp>

		<avp name="Remaining-Balance" code="2021" must="V,M" may="P" must-not="-" may-encrypt="N" vendor-id="10415">
			<data type="Grouped">
				<rule avp="Unit-Value" required="true" max="1"/>
				<rule avp="Currency-Code" required="true" max="1"/>
			</data>
		</avp>

		<avp name="AB-Response" code="7028">
			<data type="Grouped">
				<rule avp="Acct-Balance" required="true" max="1"/>
				<rule avp="Counter" required="true" max="1"/>
			</data>
		</avp>

		<avp name="Acct-Balance" code="7028">
			<data type="Grouped">
				<rule avp="Acct-Balance-Id" required="true" max="1"/>
				<rule avp="Unit-Value" required="true" max="1"/>
			</data>
		</avp>

		<avp name="Acct-Balance-Id" code="7029">
			<data type="Unsigned64"/>
		</avp>

		<avp name="Counter" code="7026">
			<data type="Grouped"/>
		</avp>
	</application>
</diameter>
	`
)
