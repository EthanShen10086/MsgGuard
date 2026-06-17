package com.ethanshen.msgguard

import android.app.Service
import android.content.Intent
import android.os.IBinder
import android.util.Log
import com.ethanshen.msgguard.network.GatewayClient
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.launch

/**
 * Stub SMS filter service. Production wiring adds an SMS receiver / default-SMS role
 * and calls [classify] before surfacing notifications.
 */
class FilterService : Service() {

    private val scope = CoroutineScope(SupervisorJob() + Dispatchers.IO)
    private val gateway by lazy { GatewayClient(BuildConfig.API_BASE) }

    override fun onBind(intent: Intent?): IBinder? = null

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        val sender = intent?.getStringExtra(EXTRA_SENDER).orEmpty()
        val body = intent?.getStringExtra(EXTRA_BODY).orEmpty()
        val deviceId = intent?.getStringExtra(EXTRA_DEVICE_ID).orEmpty()
        scope.launch {
            runCatching { classify(sender, body, deviceId) }
                .onFailure { Log.w(TAG, "classify failed", it) }
            stopSelf(startId)
        }
        return START_NOT_STICKY
    }

    suspend fun classify(sender: String, body: String, deviceId: String): GatewayClient.ClassifyResult {
        gateway.ensureDeviceToken(deviceId)
        return gateway.classify(sender, body)
    }

    companion object {
        private const val TAG = "FilterService"
        const val EXTRA_SENDER = "sender"
        const val EXTRA_BODY = "body"
        const val EXTRA_DEVICE_ID = "device_id"
    }
}
