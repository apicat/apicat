<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('doc_shares', function (Blueprint $table) {
            $table->id();
            $table->bigInteger('project_id')->index()->comment('项目id');
            $table->bigInteger('user_id')->comment('分享用户id');
            $table->bigInteger('doc_id')->comment('文档id');
            $table->string('secret_key')->comment('访问秘钥');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('doc_shares');
    }
};
