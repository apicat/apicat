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
        Schema::create('project_stars', function (Blueprint $table) {
            $table->id();
            $table->bigInteger('user_id')->index()->comment('用户id');
            $table->bigInteger('project_id')->comment('项目id');
            $table->integer('display_order')->comment('显示顺序');
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
        Schema::dropIfExists('project_stars');
    }
};
