package com.ethanshen.msgguard.network

import org.json.JSONObject
import java.net.HttpURLConnection
import java.net.URL

/**
 * Minimal gateway client: device bootstrap + classify/defer.
 */
class GatewayClient(private val baseUrl: String) {

    private var accessToken: String? = null

    data class ClassifyResult(
        val action: String,
        val category: String,
        val confidence: Double,
    )

    suspend fun ensureDeviceToken(deviceId: String) {
        if (!accessToken.isNullOrBlank()) return
        val body = JSONObject().put("device_id", deviceId).toString()
        val json = post("/api/v1/auth/device", body, auth = false)
        accessToken = json.getString("access_token")
    }

    suspend fun classify(sender: String, body: String): ClassifyResult {
        val payload = JSONObject()
            .put("sender", sender)
            .put("body", body)
            .toString()
        val json = post("/api/v1/classify/defer", payload, auth = true)
        return ClassifyResult(
            action = json.getString("action"),
            category = json.getString("category"),
            confidence = json.getDouble("confidence"),
        )
    }

    private fun post(path: String, body: String, auth: Boolean): JSONObject {
        val conn = (URL(baseUrl.trimEnd('/') + path).openConnection() as HttpURLConnection).apply {
            requestMethod = "POST"
            doOutput = true
            setRequestProperty("Content-Type", "application/json")
            if (auth) {
                val token = accessToken ?: error("device token required")
                setRequestProperty("Authorization", "Bearer $token")
            }
        }
        conn.outputStream.use { it.write(body.toByteArray(Charsets.UTF_8)) }
        val code = conn.responseCode
        val stream = if (code in 200..299) conn.inputStream else conn.errorStream
        val text = stream.bufferedReader().readText()
        if (code !in 200..299) {
            error("HTTP $code: $text")
        }
        return JSONObject(text)
    }
}
