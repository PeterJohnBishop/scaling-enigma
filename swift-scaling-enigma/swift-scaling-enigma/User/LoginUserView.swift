//
//  LoginUserView.swift
//  swift-scaling-enigma
//
//  Created by Peter Bishop on 3/20/25.
//

import Foundation
import SwiftUI

struct LoginUserView: View {
    @State var userVM: UserViewModel = UserViewModel()
    @State var password: String = ""
    @State var newUser: Bool = false
    @State var showAlert: Bool = false
    @State var loginSuccess: Bool = false
        
    var body: some View {
        NavigationStack{
                    VStack{
                        Spacer()
                        Text("Login").font(.system(size: 34))
                            .fontWeight(.ultraLight)
                        Divider().padding()
                        TextField("Email", text: $userVM.user.email)
                            .tint(.black)
                            .autocapitalization(.none)
                            .disableAutocorrection(true)
                            .padding()
                        
                        SecureField("Password", text:  $password)
                            .tint(.black)
                            .autocapitalization(.none)
                            .disableAutocorrection(true)
                            .padding()
                        
                        
                        Button("Submit", action: {
                            userVM.user.password = password
                            Task{
                                loginSuccess = await userVM.Login()
                            }
                        })
                        .fontWeight(.ultraLight)
                        .foregroundColor(.black)
                        .padding()
                        .background(
                            RoundedRectangle(cornerRadius: 8)
                                .fill(Color.white)
                                .shadow(color: .gray.opacity(0.4), radius: 4, x: 2, y: 2)
                        )
                        .onChange(of: loginSuccess, {
                            oldResponse, newResponse in
                            if !newResponse {
                                showAlert = true
                            }
                        })
                        .alert("Error", isPresented: $showAlert) {
                                        Button("OK", role: .cancel) {
                                            userVM.user.email = ""
                                            userVM.user.password = ""
                                        }
                                    } message: {
                                        Text(String(userVM.error ?? "Server Error!"))
                                    }
                                    .navigationDestination(isPresented: $loginSuccess, destination: {
                                        SuccessView().navigationBarBackButtonHidden(true)
                                    })
                        Spacer()
                        HStack{
                            Spacer()
                            Text("I don't have an account.").fontWeight(.ultraLight)
                            Button("Register", action: {
                                newUser = true
                            }).foregroundStyle(.black)
                                .fontWeight(.light)
                                .navigationDestination(isPresented: $newUser, destination: {
                                    CreateUserView().navigationBarBackButtonHidden(true)
                                })
                            Spacer()
                        }
                    }.onAppear{
                        
                    }
                }
    }
}

#Preview {
    LoginUserView()
}
