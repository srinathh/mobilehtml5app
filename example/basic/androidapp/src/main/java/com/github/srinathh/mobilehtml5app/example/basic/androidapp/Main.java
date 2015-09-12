
package com.github.srinathh.mobilehtml5app.example.basic.androidapp;

import android.app.Activity;
import android.os.Bundle;
import android.view.KeyEvent;
import android.widget.Toast;

import org.xwalk.core.XWalkView;
import go.basic.Basic;

public class Main extends Activity {
    private XWalkView mWebView;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        mWebView = new XWalkView(this, this);
        setContentView(mWebView);
    }

	// We start the server on onResume
    @Override
    protected void onResume() {
        super.onResume();
        if (mWebView != null) {
            mWebView.resumeTimers();
            mWebView.onShow();
        }

        try {
			mWebView.load(Basic.Start() + "/", null);
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
        super.onPause();
        if (mWebView != null) {
            mWebView.pauseTimers();
            mWebView.onHide();
        }
		Basic.Stop();
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        if (mWebView != null) {
            mWebView.onDestroy();
        }
    }

    // We override back key press to close the app rather than pass it to the XWalkView to give
    // a consistent user experience with how apps behave on Android.
    // Also see https://crosswalk-project.org/jira/browse/XWALK-4816
    @Override
    public boolean dispatchKeyEvent(KeyEvent event) {
        if(event.getKeyCode() == KeyEvent.KEYCODE_BACK){
            this.finish();
            return true;
        }
        return super.dispatchKeyEvent(event);
    }
}
