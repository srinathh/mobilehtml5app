
package com.github.srinathh.mobilehtml5app.example.todoapp.androidapp;

import android.app.Activity;
import android.os.Bundle;
import android.util.Log;
import android.view.KeyEvent;
import android.webkit.WebSettings;
import android.webkit.WebViewClient;
import android.widget.Toast;

import android.webkit.WebView;

import java.io.File;

import go.todoapp.Todoapp;

public class Main extends Activity {
    private WebView mWebView;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        Log.d("lifecycleLogging","in onCreate");
        super.onCreate(savedInstanceState);
        mWebView = new WebView(this);
		WebSettings webSettings = mWebView.getSettings();
		webSettings.setJavaScriptEnabled(true);
		mWebView.setWebViewClient(new WebViewClient());
        setContentView(mWebView);
    }

	// We start the server on onResume
    @Override
    protected void onResume() {
        Log.d("lifecycleLogging","in onResume");
        super.onResume();
        try {
            File fil = this.getFilesDir();
            mWebView.loadUrl(Todoapp.Start(fil.getPath()) + "/", null);
        } catch (Exception e) {
            Toast.makeText(this,"Error:"+e.toString(),Toast.LENGTH_LONG).show();
            e.printStackTrace();
            this.finish();
        }
    }

    // Send a graceful shut down signal to the server. onPause is guaranteed
	// to be called by Android while onStop or onDestroy may not be called.
    @Override
    protected void onPause() {
        Log.d("lifecycleLogging","in onPause");
        super.onPause();
		Todoapp.Stop();
    }

    @Override
    protected void onDestroy() {
        Log.d("lifecycleLogging","in onDestroy");
        super.onDestroy();
    }

    // We override back key press to close the app rather than pass it to the WebView
    @Override
    public boolean dispatchKeyEvent(KeyEvent event) {
        if(event.getKeyCode() == KeyEvent.KEYCODE_BACK){
            this.finish();
            return true;
        }
        return super.dispatchKeyEvent(event);
    }
}
